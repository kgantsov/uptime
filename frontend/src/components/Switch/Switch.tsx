import styles from './MonitorPage.module.css';

interface Props {
    value: boolean | undefined;
    onChange: (checked:boolean) => void;
}

export default function Switch({ value, onChange }: Props): JSX.Element {
  return (
    <label className={styles.switch}>
      <input
        type="checkbox"
        checked={value}
        onChange={(e) => onChange(!value)}
      />
      <span className={`${styles.slider} ${styles.round}`}></span>
    </label>
  );
}
