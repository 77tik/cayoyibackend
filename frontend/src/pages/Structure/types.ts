// 类型和接口定义
export interface SingleModelInfoProps {
  activeCard: string;
  setActiveCard: (card: string) => void;
  queryParams: { effective_head: number; active_power: number } | null;
}

export interface SimulationResult {
  mesh_json: string;
  deplace: {
    deplace: string;
    max_displacement_location: string;
    max: number;
  };
  contrainte: {
    contrainte: string;
    max_stress_location: string;
    max: number;
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
