import { describe, expect, it } from "vitest";

import type { SceneProperty } from "../Engine/ref";

import {
  applyBackwardCompatibility,
  migrateTileType,
  migrateTerrainType,
} from "./backwardCompatibility";

describe("migrateTileType", () => {
  it('should migrate "default" to "cesium_ion" with asset id 2', () => {
    const tile = {
      id: "test-tile",
      tile_type: "default",
    };

    const result = migrateTileType(tile);

    expect(result).toEqual({
      id: "test-tile",
      tile_type: "cesium_ion",
      cesium_ion_asset_id: 2,
    });
  });

  it('should migrate "default_label" to "cesium_ion" with asset id 3', () => {
    const tile = {
      id: "test-tile",
      tile_type: "default_label",
    };

    const result = migrateTileType(tile);

    expect(result).toEqual({
      id: "test-tile",
      tile_type: "cesium_ion",
      cesium_ion_asset_id: 3,
    });
  });

  it('should migrate "default_road" to "cesium_ion" with asset id 4', () => {
    const tile = {
      id: "test-tile",
      tile_type: "default_road",
    };

    const result = migrateTileType(tile);

    expect(result).toEqual({
      id: "test-tile",
      tile_type: "cesium_ion",
      cesium_ion_asset_id: 4,
    });
  });

  it('should migrate "black_marble" to "cesium_ion" with asset id 3812', () => {
    const tile = {
      id: "test-tile",
      tile_type: "black_marble",
    };

    const result = migrateTileType(tile);

    expect(result).toEqual({
      id: "test-tile",
      tile_type: "cesium_ion",
      cesium_ion_asset_id: 3812,
    });
  });

  it('should migrate "stamen_toner" to "open_street_map"', () => {
    const tile = {
      id: "test-tile",
      tile_type: "stamen_toner",
    };

    const result = migrateTileType(tile);

    expect(result).toEqual({
      id: "test-tile",
      tile_type: "open_street_map",
    });
  });

  it('should migrate "esri_world_topo" to "open_street_map"', () => {
    const tile = {
      id: "test-tile",
      tile_type: "esri_world_topo",
    };

    const result = migrateTileType(tile);

    expect(result).toEqual({
      id: "test-tile",
      tile_type: "open_street_map",
    });
  });

  it("should not override existing cesium_ion_asset_id", () => {
    const tile = {
      id: "test-tile",
      tile_type: "default",
      cesium_ion_asset_id: 999,
    };

    const result = migrateTileType(tile);

    expect(result.cesium_ion_asset_id).toBe(999);
  });

  it("should not modify unknown tile types", () => {
    const tile = {
      id: "test-tile",
      tile_type: "unknown_type",
      tile_url: "https://example.com",
    };

    const result = migrateTileType(tile);

    expect(result).toEqual(tile);
  });

  it("should not modify tile without tile_type", () => {
    const tile = {
      id: "test-tile",
      tile_url: "https://example.com",
    };

    const result = migrateTileType(tile);

    expect(result).toEqual({
      id: "test-tile",
      tile_url: "https://example.com",
    });
  });
});

describe("migrateTerrainType", () => {
  it('should migrate "arcgis" (legacy) to "reearth_terrain"', () => {
    const terrain = {
      terrain: true,
      terrainType: "arcgis" as any, // Legacy type not in current union
    };

    const result = migrateTerrainType(terrain);

    expect(result).toEqual({
      terrain: true,
      terrainType: "reearth_terrain",
    });
  });

  it("should not modify terrain when disabled", () => {
    const terrain = {
      terrain: false,
    };

    const result = migrateTerrainType(terrain);

    expect(result).toEqual(terrain);
  });

  it("should not modify terrain with existing valid terrainType", () => {
    const terrain = {
      terrain: true,
      terrainType: "cesium" as const,
      terrainExaggeration: 1.5,
    };

    const result = migrateTerrainType(terrain);

    expect(result).toEqual(terrain);
  });

  it("should preserve other terrain properties", () => {
    const terrain = {
      terrain: true,
      terrainType: "arcgis" as any, // Legacy type not in current union
      terrainExaggeration: 2.0,
      depthTestAgainstTerrain: true,
    };

    const result = migrateTerrainType(terrain);

    expect(result).toEqual({
      terrain: true,
      terrainType: "reearth_terrain",
      terrainExaggeration: 2.0,
      depthTestAgainstTerrain: true,
    });
  });
});

describe("applyBackwardCompatibility", () => {
  it("should return undefined when sceneProperty is undefined", () => {
    const result = applyBackwardCompatibility(undefined);
    expect(result).toBeUndefined();
  });

  it("should migrate tiles in sceneProperty", () => {
    const sceneProperty: SceneProperty = {
      tiles: [
        { id: "tile-1", tile_type: "default" },
        { id: "tile-2", tile_type: "stamen_toner" },
      ],
    };

    const result = applyBackwardCompatibility(sceneProperty);

    expect(result?.tiles).toEqual([
      { id: "tile-1", tile_type: "cesium_ion", cesium_ion_asset_id: 2 },
      { id: "tile-2", tile_type: "open_street_map" },
    ]);
  });

  it("should migrate terrain in sceneProperty", () => {
    const sceneProperty: SceneProperty = {
      terrain: {
        terrain: true,
        terrainType: "arcgis" as any, // Legacy type not in current union
      },
    };

    const result = applyBackwardCompatibility(sceneProperty);

    expect(result?.terrain).toEqual({
      terrain: true,
      terrainType: "reearth_terrain",
    });
  });

  it("should migrate both tiles and terrain", () => {
    const sceneProperty: SceneProperty = {
      tiles: [{ id: "tile-1", tile_type: "default" }],
      terrain: {
        terrain: true,
        terrainType: "arcgis" as any, // Legacy type not in current union
      },
    };

    const result = applyBackwardCompatibility(sceneProperty);

    expect(result?.tiles).toEqual([
      { id: "tile-1", tile_type: "cesium_ion", cesium_ion_asset_id: 2 },
    ]);
    expect(result?.terrain).toEqual({
      terrain: true,
      terrainType: "reearth_terrain",
    });
  });

  it("should not mutate original sceneProperty", () => {
    const sceneProperty: SceneProperty = {
      tiles: [{ id: "tile-1", tile_type: "default" }],
    };

    const result = applyBackwardCompatibility(sceneProperty);

    // Original should not be modified
    expect(sceneProperty.tiles?.[0].tile_type).toBe("default");
    // Result should be migrated
    expect(result?.tiles?.[0].tile_type).toBe("cesium_ion");
  });

  it("should preserve other scene properties", () => {
    const sceneProperty: SceneProperty = {
      default: {
        camera: { lat: 0, lng: 0, height: 1000, heading: 0, pitch: 0, roll: 0, fov: 1 },
        bgcolor: "#000000",
      },
      tiles: [{ id: "tile-1", tile_type: "default" }],
    };

    const result = applyBackwardCompatibility(sceneProperty);

    expect(result?.default).toEqual(sceneProperty.default);
  });
});
