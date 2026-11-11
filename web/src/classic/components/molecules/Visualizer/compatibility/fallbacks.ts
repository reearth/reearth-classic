import { cloneDeep } from "lodash-es";

import type { SceneProperty, TerrainProperty } from "../Engine/ref";
import type { Layer } from "../Layers";

/**
 * Apply fallback transformations when Cesium Ion token is not available
 * This is a temporary measure and will be removed in the future
 */
export function applyFallbacks(
  sceneProperty: SceneProperty | undefined,
): SceneProperty | undefined {
  if (!sceneProperty) return undefined;

  const hasCesiumIonToken = !!sceneProperty.default?.ion;

  // If Cesium Ion token is available, no fallback needed
  if (hasCesiumIonToken) return sceneProperty;

  const cloned = cloneDeep(sceneProperty);

  // Apply tile fallbacks
  if (cloned.tiles) {
    cloned.tiles = cloned.tiles.map(tile => fallbackTileType(tile));
  }

  // Apply terrain fallback
  if (cloned.terrain) {
    cloned.terrain = fallbackTerrainType(cloned.terrain);
  }

  return cloned;
}

/**
 * Fallback tile type when Cesium Ion token is not available
 * Fallback rules (only when no Cesium Ion token):
 * - cesium_ion asset_id 2 → google_satellite
 * - cesium_ion asset_id 3 → google_satellite
 * - cesium_ion asset_id 4 → google_roadmap
 * - cesium_ion asset_id 3812 → nasa_black_marble
 */
export function fallbackTileType(
  tile: NonNullable<SceneProperty["tiles"]>[number],
): NonNullable<SceneProperty["tiles"]>[number] {
  if (tile.tile_type !== "cesium_ion") {
    return tile;
  }

  const assetId = tile.cesium_ion_asset_id;

  // Convert to number for comparison
  const assetIdNum = typeof assetId === "string" ? parseInt(assetId, 10) : assetId;

  switch (assetIdNum) {
    case 2:
    case 3:
      console.warn(
        `[Re:Earth] Tile type fallback: "cesium_ion" (asset_id: ${assetIdNum}) → "google_satellite" - No Cesium Ion token available (tile ID: ${tile.id})`,
      );
      return {
        ...tile,
        tile_type: "google_satellite",
      };

    case 4:
      console.warn(
        `[Re:Earth] Tile type fallback: "cesium_ion" (asset_id: 4) → "google_roadmap" - No Cesium Ion token available (tile ID: ${tile.id})`,
      );
      return {
        ...tile,
        tile_type: "google_roadmap",
      };

    case 3812:
      console.warn(
        `[Re:Earth] Tile type fallback: "cesium_ion" (asset_id: 3812) → "nasa_black_marble" - No Cesium Ion token available (tile ID: ${tile.id})`,
      );
      return {
        ...tile,
        tile_type: "nasa_black_marble",
      };

    default:
      return tile;
  }
}

/**
 * Fallback terrain type when Cesium Ion token is not available
 * Fallback rule (only when no Cesium Ion token):
 * - terrainType "cesium" → "reearth_terrain"
 */
export function fallbackTerrainType(terrain: TerrainProperty): TerrainProperty {
  if (terrain.terrainType === "cesium") {
    console.warn(
      `[Re:Earth] Terrain type fallback: "cesium" → "reearth_terrain" - No Cesium Ion token available`,
    );
    return {
      ...terrain,
      terrainType: "reearth_terrain",
    };
  }

  return terrain;
}

/**
 * Fallback layer sourceType when Cesium Ion token is not available
 * Fallback rule (only when no Cesium Ion token):
 * - "osm" → "reearth-buildings"
 *
 * @param sourceType - The layer sourceType from layer property
 * @param hasCesiumIonToken - Whether a Cesium Ion token is available
 * @param layerId - Optional layer ID for logging purposes
 * @returns The fallback sourceType or the original if no fallback needed
 */
export function fallbackLayerSourceType(
  sourceType: string | undefined,
  hasCesiumIonToken: boolean,
  layerId?: string,
): string | undefined {
  // If Cesium Ion token is available, no fallback needed
  if (hasCesiumIonToken || sourceType !== "osm") {
    return sourceType;
  }

  console.warn(
    `[Re:Earth] Layer sourceType fallback: "osm" (OSM Buildings) → "reearth-buildings" (Re:Earth Buildings) - No Cesium Ion token available${
      layerId ? ` (layer ID: ${layerId})` : ""
    }`,
  );

  return "reearth-buildings";
}

/**
 * Apply fallback transformations to a layer tree when Cesium Ion token is not available
 * This recursively transforms all layers in the tree and their children
 *
 * @param rootLayer - The root layer to transform
 * @param hasCesiumIonToken - Whether a Cesium Ion token is available
 * @returns The transformed root layer with fallbacks applied, or undefined if no root layer
 */
export function applyLayerFallbacks(
  rootLayer: Layer | undefined,
  hasCesiumIonToken: boolean,
): Layer | undefined {
  if (!rootLayer) return undefined;

  // If Cesium Ion token is available, no fallback needed
  if (hasCesiumIonToken) return rootLayer;

  // Helper function to transform a single layer
  const transformLayer = (layer: Layer): Layer => {
    // Clone the layer to avoid mutation
    const transformedLayer = { ...layer };

    // Apply sourceType fallback for tileset layers
    if (layer.extensionId === "tileset" && layer.property?.default?.sourceType === "osm") {
      transformedLayer.property = {
        ...layer.property,
        default: {
          ...layer.property.default,
          sourceType: fallbackLayerSourceType(
            layer.property.default.sourceType,
            hasCesiumIonToken,
            layer.id,
          ),
        },
      };
    }

    // Recursively transform child layers
    if (layer.children && layer.children.length > 0) {
      transformedLayer.children = layer.children.map(transformLayer);
    }

    return transformedLayer;
  };

  return transformLayer(rootLayer);
}
