import { isEqual } from "lodash-es";
import { FC, useEffect, useMemo, useState } from "react";

import { useT } from "@reearth/services/i18n";
import { styled, usePublishTheme } from "@reearth/services/theme";

import { ComponentProps as WidgetProps } from "..";
import { Credits } from "../../Engine/ref";

import ContentModal from "./ContentModal";
import useCredits, { ProcessedCredit } from "./useCredits";

export type Props = WidgetProps<DataAttributionWidgetProperty>;

export type DataAttributionWidgetProperty = {
  default?: {
    id: string;
    description?: string;
    logo?: string;
    creditUrl?: string;
  }[];
};

const DataAttribution: FC<Props> = ({
  onGetCredits,
  sceneProperty,
  widget,
  hasVisibleReearthBuildingsLayers,
}) => {
  const t = useT();
  const [visible, setVisible] = useState(false);
  const [credits, setCredits] = useState<Credits>();

  const publishedTheme = usePublishTheme(sceneProperty?.theme);

  // Check for visible reearth-buildings tileset layers
  // This effect will run whenever hasVisibleReearthBuildingsLayers changes (layer added/removed/visibility changed)
  useEffect(() => {
    if (hasVisibleReearthBuildingsLayers) {
      console.log(
        "[Re:Earth] Visible reearth-buildings tileset layers detected",
        hasVisibleReearthBuildingsLayers,
      );
    }
  }, [hasVisibleReearthBuildingsLayers]);

  const layerCredits: ProcessedCredit[] = useMemo(() => {
    if (!hasVisibleReearthBuildingsLayers) return [];
    return [
      {
        builtinHtml:
          '<a href="https://buildings.reearth.land/" target="_blank" rel="noopener">Re:Earth Buildings</a> — Buildings © <a href="https://www.openstreetmap.org/copyright" target="_blank" rel="noopener">OpenStreetMap contributors</a>, <a href="https://overturemaps.org/" target="_blank" rel="noopener">Overture Maps Foundation</a> (ODbL)<br/> Terrain by <a href="https://terrain.reearth.land/" target="_blank" rel="noopener">Re:Earth Terrain</a> (Mapterhorn / EGM2008)',
      },
    ];
  }, [hasVisibleReearthBuildingsLayers]);

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

  const { cesiumCredit, otherCredits, googleCredit } = useCredits({
    credits,
    property: widget.property,
  });

  const modalCredits = useMemo(() => {
    return [...(otherCredits ?? []), ...layerCredits];
  }, [layerCredits, otherCredits]);

  return (
    <Wrapper>
      {cesiumCredit && (
        <CesiumLink target="_blank" href={cesiumCredit.creditUrl} rel="noreferrer">
          <img src={cesiumCredit.logo} title={cesiumCredit.description} />
        </CesiumLink>
      )}
      {googleCredit && (
        <GoogleLink target="_blank" rel="noreferrer" aria-label={googleCredit.description}>
          <img
            src={googleCredit.logo}
            alt={googleCredit.description}
            title={googleCredit.description}
          />
        </GoogleLink>
      )}
      <DataLink onClick={() => setVisible(true)}>{t("Data Attribution")}</DataLink>
      <ContentModal
        publishedTheme={publishedTheme}
        visible={visible}
        credits={modalCredits}
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

const GoogleLink = styled("a")(() => ({
  height: 18,
}));

export default DataAttribution;
