import {
  ImageryProvider,
  IonImageryProvider,
  OpenStreetMapImageryProvider,
  IonWorldImageryStyle,
  UrlTemplateImageryProvider,
} from "cesium";

import { config } from "@reearth/services/config";

export const tiles = {
  cesium_ion: ({ cesiumIonAssetId, cesiumIonAccessToken } = {}) => {
    if (!cesiumIonAssetId) return null;
    const NumberAssetId = parseInt(String(cesiumIonAssetId), 10);
    if (isNaN(NumberAssetId)) {
      console.warn(`Invalid cesiumIonAssetId: ${cesiumIonAssetId}`);
      return null;
    }
    return IonImageryProvider.fromAssetId(NumberAssetId, {
      accessToken: cesiumIonAccessToken,
    }).catch(err => {
      console.error(err);
      return null;
    });
  },
  google_satellite: () => {
    const customProvider = config()?.customProviders?.imagery?.providers?.find(
      p => p.id === "google_satellite",
    );
    if (!customProvider?.url) return null;

    return new UrlTemplateImageryProvider({
      url: customProvider.url,
      credit: customProvider.credit,
    });
  },
  google_roadmap: () => {
    const customProvider = config()?.customProviders?.imagery?.providers?.find(
      p => p.id === "google_roadmap",
    );
    if (!customProvider?.url) return null;

    return new UrlTemplateImageryProvider({
      url: customProvider.url,
      credit: customProvider.credit,
    });
  },
  nasa_black_marble: () => {
    const customProvider = config()?.customProviders?.imagery?.providers?.find(
      p => p.id === "nasa_black_marble",
    );
    if (!customProvider?.url) return null;

    return new UrlTemplateImageryProvider({
      url: customProvider.url,
      credit: customProvider.credit,
    });
  },
  default: ({ cesiumIonAccessToken } = {}) =>
    IonImageryProvider.fromAssetId(IonWorldImageryStyle.AERIAL, {
      accessToken: cesiumIonAccessToken,
    }).catch(console.error),
  default_label: ({ cesiumIonAccessToken } = {}) =>
    IonImageryProvider.fromAssetId(IonWorldImageryStyle.AERIAL_WITH_LABELS, {
      accessToken: cesiumIonAccessToken,
    }).catch(console.error),
  default_road: ({ cesiumIonAccessToken } = {}) =>
    IonImageryProvider.fromAssetId(IonWorldImageryStyle.ROAD, {
      accessToken: cesiumIonAccessToken,
    }).catch(console.error),
  open_street_map: () =>
    new OpenStreetMapImageryProvider({
      url: "https://tile.openstreetmap.org",
      credit: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
    }),
  black_marble: ({ cesiumIonAccessToken } = {}) =>
    IonImageryProvider.fromAssetId(3812, { accessToken: cesiumIonAccessToken }).catch(
      console.error,
    ),
  japan_gsi_standard: () =>
    new OpenStreetMapImageryProvider({
      url: "https://cyberjapandata.gsi.go.jp/xyz/std/",
      credit:
        '<a href="https://maps.gsi.go.jp/development/ichiran.html">国土地理院</a>, Shoreline data is derived from: United States. National Imagery and Mapping Agency. "Vector Map Level 0 (VMAP0)." Bethesda, MD: Denver, CO: The Agency; USGS Information Services, 1997.',
    }),
  url: ({ url } = {}) => (url ? new UrlTemplateImageryProvider({ url }) : null),
} as {
  [key: string]: (opts?: {
    url?: string;
    cesiumIonAssetId?: string | number;
    cesiumIonAccessToken?: string;
  }) => Promise<ImageryProvider> | ImageryProvider | null;
};
