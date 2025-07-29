import { useCallback, useRef } from 'react';

interface CacheItem<T> {
  data: T;
  timestamp: number;
  key: string;
}

/**
 * 数据缓存 Hook
 * @param ttl 缓存时间（毫秒），默认 5 分钟
 */
export function useDataCache<T>(ttl: number = 5 * 60 * 1000) {
  const cacheRef = useRef<Map<string, CacheItem<T>>>(new Map());

  const generateKey = useCallback((params: any[]): string => {
    return JSON.stringify(params);
  }, []);

  const getCache = useCallback(
    (key: string): T | null => {
      const cache = cacheRef.current;
      const item = cache.get(key);

      if (!item) return null;

      // 检查是否过期
      if (Date.now() - item.timestamp > ttl) {
        cache.delete(key);
        return null;
      }

      return item.data;
    },
    [ttl],
  );

  const setCache = useCallback(
    (key: string, data: T): void => {
      const cache = cacheRef.current;
      cache.set(key, {
        data,
        timestamp: Date.now(),
        key,
      });

      // 清理过期缓存（简单的清理策略）
      if (cache.size > 50) {
        const now = Date.now();
        for (const [k, item] of cache.entries()) {
          if (now - item.timestamp > ttl) {
            cache.delete(k);
          }
        }
      }
    },
    [ttl],
  );

  const getCachedData = useCallback(
    (params: any[], fetcher: () => Promise<T>): Promise<T> => {
      const key = generateKey(params);
      const cached = getCache(key);

      if (cached !== null) {
        return Promise.resolve(cached);
      }

      return fetcher().then((data) => {
        setCache(key, data);
        return data;
      });
    },
    [generateKey, getCache, setCache],
  );

  const clearCache = useCallback(() => {
    cacheRef.current.clear();
  }, []);

  return {
    getCachedData,
    clearCache,
  };
}

export default useDataCache;
