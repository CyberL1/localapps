export interface Props {
  id?: string;
  data: object;
  items: Item[];
}

export interface Item {
  title: string;
  onclick?: () => void;
}
