import { cloneDeep } from "lodash-es";

import type { SceneProperty, TerrainProperty } from "../Engine/ref";

/**
 * Apply backward compatibility transformations to scene property
 * This ensures old tile/terrain types are migrated to new provider names
 */
export function applyBackwardCompatibility(
  sceneProperty: SceneProperty | undefined,
): SceneProperty | undefined {
  if (!sceneProperty) return undefined;

  const cloned = cloneDeep(sceneProperty);

  // Apply tile type backward compatibility
  if (cloned.tiles) {
    cloned.tiles = cloned.tiles.map(tile => migrateTileType(tile));
  }

  // Apply terrain type backward compatibility
  if (cloned.terrain) {
    cloned.terrain = migrateTerrainType(cloned.terrain);
  }

  return cloned;
}

/**
 * Migrate tile type from old format to new format
 * Backward compatibility rules:
 * - "default" → "cesium_ion" with cesiumIonAssetId: 2
 * - "default_label" → "cesium_ion" with cesiumIonAssetId: 3
 * - "default_road" → "cesium_ion" with cesiumIonAssetId: 4
 * - "black_marble" → "cesium_ion" with cesiumIonAssetId: 3812
 * - "stamen_toner" → "carto_light"
 * - "esri_world_topo" → "open_street_map"
 */
export function migrateTileType(
  tile: NonNullable<SceneProperty["tiles"]>[number],
): NonNullable<SceneProperty["tiles"]>[number] {
  const tileType = tile.tile_type;

  // Backward compatibility: migrate old tile types to new ones
  switch (tileType) {
    case "default": {
      const assetId = tile.cesium_ion_asset_id ?? 2;
      console.warn(
        `[Re:Earth] Tile type migrated: "default" → "cesium_ion" (asset_id: ${assetId}) - Backward compatibility (tile ID: ${tile.id})`,
      );
      return {
        ...tile,
        tile_type: "cesium_ion",
        cesium_ion_asset_id: assetId,
      };
    }

    case "default_label": {
      const assetId = tile.cesium_ion_asset_id ?? 3;
      console.warn(
        `[Re:Earth] Tile type migrated: "default_label" → "cesium_ion" (asset_id: ${assetId}) - Backward compatibility (tile ID: ${tile.id})`,
      );
      return {
        ...tile,
        tile_type: "cesium_ion",
        cesium_ion_asset_id: assetId,
      };
    }

    case "default_road": {
      const assetId = tile.cesium_ion_asset_id ?? 4;
      console.warn(
        `[Re:Earth] Tile type migrated: "default_road" → "cesium_ion" (asset_id: ${assetId}) - Backward compatibility (tile ID: ${tile.id})`,
      );
      return {
        ...tile,
        tile_type: "cesium_ion",
        cesium_ion_asset_id: assetId,
      };
    }

    case "black_marble": {
      const assetId = tile.cesium_ion_asset_id ?? 3812;
      console.warn(
        `[Re:Earth] Tile type migrated: "black_marble" → "cesium_ion" (asset_id: ${assetId}) - Backward compatibility (tile ID: ${tile.id})`,
      );
      return {
        ...tile,
        tile_type: "cesium_ion",
        cesium_ion_asset_id: assetId,
      };
    }

    case "stamen_toner":
      console.warn(
        `[Re:Earth] Tile type migrated: "stamen_toner" → "carto_light" - Backward compatibility (tile ID: ${tile.id})`,
      );
      return {
        ...tile,
        tile_type: "carto_light",
      };

    case "esri_world_topo":
      console.warn(
        `[Re:Earth] Tile type migrated: "esri_world_topo" → "open_street_map" - Backward compatibility (tile ID: ${tile.id})`,
      );
      return {
        ...tile,
        tile_type: "open_street_map",
      };

    default:
      return tile;
  }
}

/**
 * Migrate terrain type from old format to new format
 * Backward compatibility rules:
 * - If terrainType is "arcgis" (legacy) → change to "reearth_terrain"
 */
export function migrateTerrainType(terrain: TerrainProperty): TerrainProperty {
  const { terrainType } = terrain;

  // Migrate "arcgis" to "reearth_terrain" (handles legacy string type)
  if ((terrainType as string) === "arcgis") {
    console.warn(
      `[Re:Earth] Terrain type migrated: "arcgis" → "reearth_terrain" - Backward compatibility (legacy provider)`,
    );
    return {
      ...terrain,
      terrainType: "reearth_terrain",
    };
  }

  return terrain;
}
