// 类型和接口定义
export interface SingleModelInfoProps {
  activeCard: string;
  setActiveCard: (card: string) => void;
  queryParams: { effective_head: number; active_power: number } | null;
}

export interface SimulationResult {
  data: {
    mesh_json: string;
    velocity: {
      h: string;
      v: string;
      stream_line: string;
      volute_average_velocity: number;
      max_velocity?: number;
    };
    pressure: {
      h: string;
      v: string;
      volute_pressure: number;
      max_pressure?: number;
    };
    vof: {
      mesh_json: string;
      vof: string;
      runner_cavitation_bubble_count: number;
      blade_cavitation_area: number;
    };
    vortex: {
      vortex: string;
      vortex_concentration_location: string;
    };
  };
}

export interface CardConfigItem {
  label: string;
  value: string;
}

export interface CardConfig {
  [key: string]: CardConfigItem[];
}

export interface ButtonConfig {
  key: string;
  icon: string;
  activeIcon: string;
  text: string;
}

export type ProfileType = 'h' | 'v';
