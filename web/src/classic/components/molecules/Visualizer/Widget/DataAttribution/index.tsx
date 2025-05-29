import { isEqual } from "lodash-es";
import { FC, useEffect, useState } from "react";

import { useT } from "@reearth/services/i18n";
import { styled, usePublishTheme } from "@reearth/services/theme";

import { ComponentProps as WidgetProps } from "..";
import { Credits } from "../../Engine/ref";

import ContentModal from "./ContentModal";
import useCredits from "./useCredits";

export type Props = WidgetProps<DataAttributionWidgetProperty>;

export type DataAttributionWidgetProperty = {
  default?: {
    id: string;
    description?: string;
    logo?: string;
    creditUrl?: string;
  }[];
};

const DataAttribution: FC<Props> = ({ onGetCredits, sceneProperty, widget }) => {
  const t = useT();
  const [visible, setVisible] = useState(false);
  const [credits, setCredits] = useState<Credits>();

  const publishedTheme = usePublishTheme(sceneProperty?.theme);

  useEffect(() => {
    let intervalId: NodeJS.Timeout | null = null;

    const newCredits = onGetCredits?.();
    if (newCredits && !isEqual(newCredits, credits)) {
      setCredits(newCredits);
    }

    intervalId = setInterval(() => {
      const newCredits = onGetCredits?.();
      if (newCredits && !isEqual(newCredits, credits)) {
        setCredits(newCredits);
      }
    }, 3000);

    return () => {
      if (intervalId) clearInterval(intervalId);
    };
  }, [onGetCredits, credits]);

  const { cesiumCredit, otherCredits } = useCredits({ credits, property: widget.property });

  return (
    <Wrapper>
      {cesiumCredit && (
        <CesiumLink target="_blank" href={cesiumCredit.creditUrl} rel="noreferrer">
          <img src={cesiumCredit.logo} title={cesiumCredit.description} />
        </CesiumLink>
      )}
      <DataLink onClick={() => setVisible(true)}>{t("Data Attribution")}</DataLink>
      <ContentModal
        publishedTheme={publishedTheme}
        visible={visible}
        credits={otherCredits}
        onClose={() => setVisible(false)}
      />
    </Wrapper>
  );
};

const CesiumLink = styled("a")(() => ({
  height: 28,
  display: "block",
  padding: 4,
}));

const Wrapper = styled("div")(({ theme }) => ({
  display: "flex",
  alignItems: "center",
  gap: theme.spacing.small,
  fontSize: 12,
  fontWeight: "bold",
}));

const DataLink = styled("div")(() => ({
  cursor: "pointer",
  padding: 4,
}));

export default DataAttribution;
