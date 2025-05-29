import { FC } from "react";
import { createPortal } from "react-dom";

import Icon from "@reearth/classic/components/atoms/Icon";
import { PublishTheme } from "@reearth/classic/theme";
import { useT } from "@reearth/services/i18n";
import { styled } from "@reearth/services/theme";

import { VisualizerDOMId } from "../..";

import { ProcessedCredit } from "./useCredits";

type ContentModalProps = {
  visible?: boolean;
  credits?: ProcessedCredit[];
  publishedTheme?: PublishTheme;
  onClose?: () => void;
};

const ContentModal: FC<ContentModalProps> = ({ visible, credits, publishedTheme, onClose }) => {
  const portalRoot = document.getElementById(VisualizerDOMId);
  const t = useT();

  return visible && portalRoot
    ? createPortal(
        <Wrapper>
          <Panel publishedTheme={publishedTheme}>
            <Header>
              <Title>{t("Data Provided by:")}</Title>
              <CloseButton>
                <Icon icon="cancel" size={24} onClick={onClose} />
              </CloseButton>
            </Header>

            <Content>
              {credits?.map((credit, index) => (
                <ListItem key={index}>
                  <ListMarker>•</ListMarker>
                  {credit.creditUrl ? (
                    <CreditItemLink
                      target="_blank"
                      href={credit.creditUrl}
                      rel="noopener noreferrer">
                      {credit.logo && (
                        <LogoWrapper>
                          <StyledImage src={credit.logo} />
                        </LogoWrapper>
                      )}
                      {credit.description && <CreditText>{credit.description}</CreditText>}
                    </CreditItemLink>
                  ) : (
                    <CreditItem>
                      {credit.logo && (
                        <LogoWrapper>
                          <StyledImage src={credit.logo} />
                        </LogoWrapper>
                      )}
                      {credit.description && <CreditText>{credit.description}</CreditText>}
                    </CreditItem>
                  )}
                </ListItem>
              ))}
            </Content>
          </Panel>
        </Wrapper>,
        portalRoot,
      )
    : null;
};

const Wrapper = styled("div")(({ theme }) => ({
  position: "absolute",
  width: "100%",
  height: "100%",
  left: 0,
  top: 0,
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  backgroundColor: "rgba(0, 0, 0, 0.7)",
  zIndex: theme.classic.zIndexes.dataAttribution,
}));

const Panel = styled("div")<{ publishedTheme?: PublishTheme }>(({ theme, publishedTheme }) => ({
  backgroundColor: publishedTheme?.background || theme.classic.colors.dark.bg[1],
  color: publishedTheme?.mainText || theme.classic.colors.dark.text.main,
  fontSize: 12,
  borderRadius: 8,
  boxShadow: "0 4px 6px rgba(0, 0, 0, 0.1)",
  width: "100%",
  maxWidth: 400,
  display: "flex",
  flexDirection: "column",
  gap: theme.spacing.normal,
  "& a": {
    color: publishedTheme?.mainText || theme.classic.colors.dark.text.main,
  },
}));

const Header = styled("div")(() => ({
  display: "flex",
  alignItems: "center",
  justifyContent: "space-between",
  position: "relative",
  padding: "20px 20px 0 20px",
}));

const Title = styled("div")(() => ({
  fontSize: 14,
  fontWeight: "bold",
}));

const CloseButton = styled("div")(() => ({
  cursor: "pointer",
}));

const Content = styled("div")(() => ({
  display: "flex",
  flexDirection: "column",
  gap: 8,
  height: "100%",
  maxHeight: "60vh",
  overflowY: "auto",
  padding: "0 20px 20px 20px",
}));

const ListItem = styled("li")(() => ({
  display: "flex",
  alignItems: "center",
  gap: 12,
}));

const ListMarker = styled("div")(() => ({
  width: 5,
}));

const CreditItem = styled("div")(() => ({
  display: "flex",
  alignItems: "center",
  justifyContent: "flex-start",
  gap: 8,
}));

const CreditItemLink = styled("a")(() => ({
  display: "flex",
  alignItems: "center",
  justifyContent: "flex-start",
  gap: 8,
}));

const CreditText = styled("span")(() => ({
  display: "block",
}));

const LogoWrapper = styled("div")(() => ({
  position: "relative",
  boxSizing: "border-box",
  padding: 4,
  backgroundColor: "rgba(0, 0, 0, 0.2)",
  borderRadius: 4,
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
}));

const StyledImage = styled("img")(() => ({
  maxHeight: 30,
  width: "auto",
}));

export default ContentModal;
