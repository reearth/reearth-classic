import { FC } from "react";

import { styled } from "@reearth/services/theme";

import { IconName, Icon } from "../Icon";

export type ButtonProps = {
  appearance?: "primary" | "secondary" | "dangerous" | "simple";
  disabled?: boolean;
  size?: "normal" | "small";
  iconButton?: boolean;
  icon?: IconName;
  title?: string;
  extendWidth?: boolean;
  minWidth?: number;
  onClick?: () => void;
};

export const Button: FC<ButtonProps> = ({
  appearance = "secondary",
  icon,
  disabled,
  size = "normal",
  iconButton,
  title,
  extendWidth,
  minWidth,
  onClick,
}) => {
  return (
    <StyledButton
      disabled={disabled}
      appearance={appearance}
      size={size}
      iconButton={iconButton}
      extendWidth={extendWidth}
      minWidth={minWidth}
      onClick={onClick}>
      {icon && <Icon icon={icon} />}
      {!iconButton && title}
    </StyledButton>
  );
};

const StyledButton = styled("button")<{
  size: "normal" | "small";
  appearance: "primary" | "secondary" | "dangerous" | "simple";
  iconButton?: boolean;
  extendWidth?: boolean;
  minWidth?: number;
}>(({ appearance, size, iconButton, extendWidth, minWidth, theme }) => ({
  display: "flex",
  flexDirection: "row",
  alignItems: "center",
  justifyContent: "center",
  gap: `${theme.spacing.small}px`,
  border:
    appearance === "simple"
      ? "none"
      : `1px solid ${
          appearance === "secondary"
            ? `${theme.outline.weak}`
            : appearance === "primary"
            ? `${theme.primary.main}`
            : `${theme.dangerous.main}`
        }`,
  padding:
    size === "small"
      ? iconButton
        ? `${theme.spacing.micro}px`
        : `${theme.spacing.smallest}px ${theme.spacing.small}px`
      : iconButton
      ? `${theme.spacing.small}px`
      : `${theme.spacing.small}px ${theme.spacing.large}px`,
  borderRadius: size === "small" ? `${theme.radius.small}px` : `${theme.radius.normal}px`,
  color:
    appearance === "primary"
      ? `${theme.primary.main}`
      : appearance === "dangerous"
      ? `${theme.dangerous.main}`
      : `${theme.content.main}`,
  backgroundColor: appearance === "simple" ? "transparent" : `${theme.bg[1]}`,
  width: !extendWidth ? "fit-content" : "100%",
  minWidth: minWidth ? `${minWidth}px` : "",
  boxShadow: theme.shadow.input,
  ["&:hover"]: {
    borderColor: "transparent",
    color: `${theme.content.withBackground}`,
    backgroundColor: appearance === "simple" ? "transparent" : `${theme[appearance].weak}`,
  },
  ["&:active"]: {
    borderColor: "transparent",
    color: `${theme.content.withBackground}`,
    backgroundColor: appearance === "simple" ? "transparent" : `${theme[appearance].main}`,
    boxShadow: "none",
  },
  ["&:disabled"]: {
    cursor: "not-allowed",
    borderColor: "transparent",
    color: `${theme.content.weaker}`,
    backgroundColor: appearance !== "simple" ? `${theme.bg[1]}` : "transparent",
    boxShadow: "none",
  },
}));
