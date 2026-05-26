import { ImageryProvider } from "cesium";
import { isEqual } from "lodash-es";
import { useCallback, useMemo, useRef, useLayoutEffect, useState } from "react";
import { ImageryLayer } from "resium";

import { tiles as tilePresets } from "./presets";

export type ImageryLayerData = {
  id: string;
  provider: ImageryProvider;
  min?: number;
  max?: number;
  opacity?: number;
};

export type Tile = {
  id: string;
  tile_url?: string;
  tile_type?: string;
  tile_opacity?: number;
  tile_minLevel?: number;
  tile_maxLevel?: number;
  cesium_ion_asset_id?: string | number;
};

export type Props = {
  tiles?: Tile[];
  cesiumIonAccessToken?: string;
};

export default function ImageryLayers({ tiles, cesiumIonAccessToken }: Props) {
  const { providers, updated } = useImageryProviders({
    tiles,
    cesiumIonAccessToken,
    presets: tilePresets,
  });

  // force rerendering all layers when any provider is updated
  // since Resium does not sort layers according to ImageryLayer component order
  const [counter, setCounter] = useState(0);
  useLayoutEffect(() => {
    if (updated) setCounter(c => c + 1);
  }, [providers, updated]);

  return (
    <>
      {tiles
        ?.map(({ id, ...tile }) => ({ ...tile, id, provider: providers[id]?.[3] }))
        .map(({ id, tile_opacity: opacity, tile_minLevel: min, tile_maxLevel: max, provider }, i) =>
          provider ? (
            <ImageryLayer
              key={`${id}_${i}_${counter}`}
              imageryProvider={provider}
              minimumTerrainLevel={min}
              maximumTerrainLevel={max}
              alpha={opacity}
            />
          ) : null,
        )}
    </>
  );
}

type Providers = {
  [id: string]: [
    string | undefined,
    string | undefined,
    string | number | undefined,
    ImageryProvider,
  ];
};

export function useImageryProviders({
  tiles = [],
  cesiumIonAccessToken,
  presets,
}: {
  tiles?: Tile[];
  cesiumIonAccessToken?: string;
  presets: {
    [key: string]: (opts?: {
      url?: string;
      cesiumIonAccessToken?: string;
      cesiumIonAssetId?: string | number;
    }) => Promise<ImageryProvider> | ImageryProvider | null;
  };
}): { providers: Providers; updated: boolean } {
  const newTile = useCallback(
    (t: Tile, ciat?: string) =>
      presets[t.tile_type || "google_satellite"]({
        url: t.tile_url,
        cesiumIonAccessToken: ciat,
        cesiumIonAssetId: t.cesium_ion_asset_id,
      }),
    [presets],
  );

  const prevCesiumIonAccessToken = useRef(cesiumIonAccessToken);
  const tileKeys = tiles.map(t => t.id).join(",");
  const prevTileKeys = useRef(tileKeys);
  const prevProviders = useRef<Providers>({});

  // Manage TileProviders so that TileProvider does not need to be recreated each time tiles are updated.
  const { providers, updated } = useMemo(() => {
    const isCesiumAccessTokenUpdated = prevCesiumIonAccessToken.current !== cesiumIonAccessToken;
    const prevProvidersKeys = Object.keys(prevProviders.current);
    const added = tiles.map(t => t.id).filter(t => t && !prevProvidersKeys.includes(t));

    const rawProviders = [
      ...Object.entries(prevProviders.current),
      ...added.map(a => [a, undefined] as const),
    ].map(([k, v]) => ({
      key: k,
      added: added.includes(k),
      prevType: v?.[0],
      prevUrl: v?.[1],
      prevAssetId: v?.[2],
      prevProvider: v?.[3],
      tile: tiles.find(t => t.id === k),
    }));

    const providers = Object.fromEntries(
      rawProviders
        .map(
          ({
            key,
            added,
            prevType,
            prevUrl,
            prevAssetId,
            prevProvider,
            tile,
          }):
            | [
                string,
                [
                  string | undefined,
                  string | undefined,
                  string | number | undefined,
                  Promise<ImageryProvider> | ImageryProvider | null | undefined,
                ],
              ]
            | null =>
            !tile
              ? null
              : [
                  key,
                  added ||
                  prevType !== tile.tile_type ||
                  prevUrl !== tile.tile_url ||
                  prevAssetId !== tile.cesium_ion_asset_id ||
                  (isCesiumAccessTokenUpdated &&
                    // Recreate provider when token changes for tile types that use Cesium Ion
                    (!tile.tile_type ||
                      tile.tile_type === "google_satellite" ||
                      tile.tile_type === "cesium_ion" ||
                      tile.tile_type === "default" ||
                      tile.tile_type === "default_label" ||
                      tile.tile_type === "default_road" ||
                      tile.tile_type === "black_marble"))
                    ? [
                        tile.tile_type,
                        tile.tile_url,
                        tile.cesium_ion_asset_id,
                        newTile(tile, cesiumIonAccessToken),
                      ]
                    : [prevType, prevUrl, prevAssetId, prevProvider],
                ],
        )
        .filter(
          (
            e,
          ): e is [
            string,
            [string | undefined, string | undefined, string | number | undefined, ImageryProvider],
          ] => !!e?.[1][3],
        ),
    );

    const updated =
      !!added.length ||
      !!isCesiumAccessTokenUpdated ||
      !isEqual(prevTileKeys.current, tileKeys) ||
      rawProviders.some(
        p =>
          p.tile &&
          (p.prevType !== p.tile.tile_type ||
            p.prevUrl !== p.tile.tile_url ||
            p.prevAssetId !== p.tile.cesium_ion_asset_id),
      );

    prevTileKeys.current = tileKeys;
    prevCesiumIonAccessToken.current = cesiumIonAccessToken;

    return { providers, updated };
  }, [cesiumIonAccessToken, tiles, tileKeys, newTile]);

  prevProviders.current = providers;
  return { providers, updated };
}
