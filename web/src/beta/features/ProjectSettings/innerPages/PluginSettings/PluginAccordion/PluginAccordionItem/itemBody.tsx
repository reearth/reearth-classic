import Markdown from "@reearth/beta/components/Markdown";
import { useTheme, styled } from "@reearth/services/theme";

export type Props = {
  children?: string;
};

const PluginAccordionItemBody: React.FC<Props> = ({ children }) => {
  const theme = useTheme();
  const markdownTypographyStyle = {
    color: theme.classic.text.default,
  };
  return (
    <Wrapper>
      <Markdown backgroundColor={theme.classic.pluginList.bg} styles={markdownTypographyStyle}>
        {children}
      </Markdown>
    </Wrapper>
  );
};

export default PluginAccordionItemBody;

const Wrapper = styled.div`
  padding: ${({ theme }) => theme.spacing.super}px;
`;
