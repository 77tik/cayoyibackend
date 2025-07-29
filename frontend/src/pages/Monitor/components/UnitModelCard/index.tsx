import React, { useMemo } from 'react';
import { StrainPoint } from '../../index';

interface UnitModelCardProps {
  unit: {
    id: number;
    title: string;
  };
  strainData?: StrainPoint;
}

const UnitModelCard: React.FC<UnitModelCardProps> = ({ unit, strainData }) => {
  // 使用 useMemo 缓存字段键名计算
  const fieldKeys = useMemo(() => {
    const unitNo = unit.id;
    const coverKey = ['one_upper', 'two_cover', 'three_cover', 'four_cover'][
      unitNo - 1
    ] as keyof StrainPoint;
    const doorKey = ['one_door', 'two_door', 'three_door', 'four_door'][
      unitNo - 1
    ] as keyof StrainPoint;

    return { coverKey, doorKey };
  }, [unit.id]);

  // 使用 useMemo 缓存格式化后的数值
  const formattedValues = useMemo(() => {
    const { coverKey, doorKey } = fieldKeys;

    // 格式化显示数值，保留两位小数
    const formatValue = (val: number | null | undefined) =>
      val === null || val === undefined ? '--' : val.toFixed(2);

    const coverValue = strainData?.[coverKey];
    const doorValue = strainData?.[doorKey];

    return {
      cover:
        coverValue !== null && coverValue !== undefined
          ? `${formatValue(coverValue)} MW`
          : '-\n-\n-\n-',
      door:
        doorValue !== null && doorValue !== undefined
          ? `${formatValue(doorValue)} MW`
          : '-\n-\n-\n-',
    };
  }, [strainData, fieldKeys]);

  return (
    <div className="section-item">
      <div className="section-item-title">{unit.title}</div>
      <div className="section-item-content">
        <div className="cover-strain">顶盖应变：{formattedValues.cover}</div>
        <div className="strain-model">
          <img src="/images/位图.png" alt="strain-model" />
        </div>
        <div className="volute-manhole-door-strain">
          涡壳入孔门应变：{formattedValues.door}
        </div>
      </div>
    </div>
  );
};

export default React.memo(UnitModelCard);
