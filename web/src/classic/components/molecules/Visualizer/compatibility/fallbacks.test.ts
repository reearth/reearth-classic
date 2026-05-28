import { describe, expect, it } from "vitest";

import type { SceneProperty } from "../Engine/ref";
import type { Layer } from "../Layers";

import {
  applyFallbacks,
  fallbackTileType,
  fallbackTerrainType,
  fallbackLayerSourceType,
  applyLayerFallbacks,
} from "./fallbacks";

describe("fallbackTileType", () => {
  it("should fallback cesium_ion asset_id 2 to google_satellite", () => {
    const tile = {
      id: "test-tile",
      tile_type: "cesium_ion",
      cesium_ion_asset_id: 2,
    };

    const result = fallbackTileType(tile);

    expect(result).toEqual({
      id: "test-tile",
      tile_type: "google_satellite",
      cesium_ion_asset_id: 2,
    });
  });

  it("should fallback cesium_ion asset_id 3 to google_satellite", () => {
    const tile = {
      id: "test-tile",
      tile_type: "cesium_ion",
      cesium_ion_asset_id: 3,
    };

    const result = fallbackTileType(tile);

    expect(result).toEqual({
      id: "test-tile",
      tile_type: "google_satellite",
      cesium_ion_asset_id: 3,
    });
  });

  it("should fallback cesium_ion asset_id 4 to google_roadmap", () => {
    const tile = {
      id: "test-tile",
      tile_type: "cesium_ion",
      cesium_ion_asset_id: 4,
    };

    const result = fallbackTileType(tile);

    expect(result).toEqual({
      id: "test-tile",
      tile_type: "google_roadmap",
      cesium_ion_asset_id: 4,
    });
  });

  it("should fallback cesium_ion asset_id 3812 to nasa_black_marble", () => {
    const tile = {
      id: "test-tile",
      tile_type: "cesium_ion",
      cesium_ion_asset_id: 3812,
    };

    const result = fallbackTileType(tile);

    expect(result).toEqual({
      id: "test-tile",
      tile_type: "nasa_black_marble",
      cesium_ion_asset_id: 3812,
    });
  });

  it("should handle string asset_id", () => {
    const tile = {
      id: "test-tile",
      tile_type: "cesium_ion",
      cesium_ion_asset_id: "2",
    };

    const result = fallbackTileType(tile);

    expect(result.tile_type).toBe("google_satellite");
  });

  it("should not fallback cesium_ion with unknown asset_id", () => {
    const tile = {
      id: "test-tile",
      tile_type: "cesium_ion",
      cesium_ion_asset_id: 9999,
    };

    const result = fallbackTileType(tile);

    expect(result).toEqual(tile);
  });

  it("should not fallback non-cesium_ion tiles", () => {
    const tile = {
      id: "test-tile",
      tile_type: "open_street_map",
    };

    const result = fallbackTileType(tile);

    expect(result).toEqual(tile);
  });

  it("should not fallback cesium_ion without asset_id", () => {
    const tile = {
      id: "test-tile",
      tile_type: "cesium_ion",
    };

    const result = fallbackTileType(tile);

    expect(result).toEqual(tile);
  });
});

describe("fallbackTerrainType", () => {
  it('should fallback terrainType "cesium" to "reearth_terrain"', () => {
    const terrain = {
      terrain: true,
      terrainType: "cesium" as const,
    };

    const result = fallbackTerrainType(terrain);

    expect(result).toEqual({
      terrain: true,
      terrainType: "reearth_terrain",
    });
  });

  it("should not fallback other terrain types", () => {
    const terrain = {
      terrain: true,
      terrainType: "reearth_terrain" as const,
    };

    const result = fallbackTerrainType(terrain);

    expect(result).toEqual(terrain);
  });

  it("should preserve other terrain properties", () => {
    const terrain = {
      terrain: true,
      terrainType: "cesium" as const,
      terrainExaggeration: 2.0,
      depthTestAgainstTerrain: true,
    };

    const result = fallbackTerrainType(terrain);

    expect(result).toEqual({
      terrain: true,
      terrainType: "reearth_terrain",
      terrainExaggeration: 2.0,
      depthTestAgainstTerrain: true,
    });
  });

  it("should not fallback when terrain is disabled", () => {
    const terrain = {
      terrain: false,
      terrainType: "cesium" as const,
    };

    const result = fallbackTerrainType(terrain);

    expect(result).toEqual({
      terrain: false,
      terrainType: "reearth_terrain",
    });
  });
});

describe("fallbackLayerSourceType", () => {
  it('should fallback "osm" to "reearth-buildings" when no Cesium Ion token', () => {
    const result = fallbackLayerSourceType("osm", false);

    expect(result).toBe("reearth-buildings");
  });

  it('should NOT fallback "osm" when Cesium Ion token is provided', () => {
    const result = fallbackLayerSourceType("osm", true);

    expect(result).toBe("osm");
  });

  it('should not fallback non-"osm" sourceTypes', () => {
    const sourceTypes = ["url", "google-photorealistic", "reearth-buildings", undefined];

    sourceTypes.forEach(sourceType => {
      const result = fallbackLayerSourceType(sourceType, false);
      expect(result).toBe(sourceType);
    });
  });

  it("should not fallback when sourceType is undefined", () => {
    const result = fallbackLayerSourceType(undefined, false);

    expect(result).toBeUndefined();
  });

  it("should handle layerId parameter for logging", () => {
    const result = fallbackLayerSourceType("osm", false, "test-layer-id");

    expect(result).toBe("reearth-buildings");
  });

  it("should preserve sourceType when token is available even with layerId", () => {
    const result = fallbackLayerSourceType("osm", true, "test-layer-id");

    expect(result).toBe("osm");
  });
});

describe("applyFallbacks", () => {
  it("should return undefined when sceneProperty is undefined", () => {
    const result = applyFallbacks(undefined);
    expect(result).toBeUndefined();
  });

  it("should NOT apply fallbacks when Cesium Ion token is provided", () => {
    const sceneProperty: SceneProperty = {
      default: {
        ion: "my-cesium-ion-token",
      },
      tiles: [{ id: "tile-1", tile_type: "cesium_ion", cesium_ion_asset_id: 2 }],
      terrain: {
        terrain: true,
        terrainType: "cesium",
      },
    };

    const result = applyFallbacks(sceneProperty);

    // Should return the same object without fallbacks
    expect(result?.tiles?.[0].tile_type).toBe("cesium_ion");
    expect(result?.terrain?.terrainType).toBe("cesium");
  });

  it("should apply tile fallbacks when no Cesium Ion token", () => {
    const sceneProperty: SceneProperty = {
      tiles: [
        { id: "tile-1", tile_type: "cesium_ion", cesium_ion_asset_id: 2 },
        { id: "tile-2", tile_type: "cesium_ion", cesium_ion_asset_id: 4 },
      ],
    };

    const result = applyFallbacks(sceneProperty);

    expect(result?.tiles).toEqual([
      { id: "tile-1", tile_type: "google_satellite", cesium_ion_asset_id: 2 },
      { id: "tile-2", tile_type: "google_roadmap", cesium_ion_asset_id: 4 },
    ]);
  });

  it("should apply terrain fallback when no Cesium Ion token", () => {
    const sceneProperty: SceneProperty = {
      terrain: {
        terrain: true,
        terrainType: "cesium",
      },
    };

    const result = applyFallbacks(sceneProperty);

    expect(result?.terrain).toEqual({
      terrain: true,
      terrainType: "reearth_terrain",
    });
  });

  it("should apply both tile and terrain fallbacks when no Cesium Ion token", () => {
    const sceneProperty: SceneProperty = {
      tiles: [{ id: "tile-1", tile_type: "cesium_ion", cesium_ion_asset_id: 3812 }],
      terrain: {
        terrain: true,
        terrainType: "cesium",
      },
    };

    const result = applyFallbacks(sceneProperty);

    expect(result?.tiles).toEqual([
      { id: "tile-1", tile_type: "nasa_black_marble", cesium_ion_asset_id: 3812 },
    ]);
    expect(result?.terrain).toEqual({
      terrain: true,
      terrainType: "reearth_terrain",
    });
  });

  it("should not mutate original sceneProperty", () => {
    const sceneProperty: SceneProperty = {
      tiles: [{ id: "tile-1", tile_type: "cesium_ion", cesium_ion_asset_id: 2 }],
    };

    const result = applyFallbacks(sceneProperty);

    // Original should not be modified
    expect(sceneProperty.tiles?.[0].tile_type).toBe("cesium_ion");
    // Result should be modified
    expect(result?.tiles?.[0].tile_type).toBe("google_satellite");
  });

  it("should preserve other scene properties", () => {
    const sceneProperty: SceneProperty = {
      default: {
        camera: { lat: 0, lng: 0, height: 1000, heading: 0, pitch: 0, roll: 0, fov: 1 },
        bgcolor: "#000000",
      },
      tiles: [{ id: "tile-1", tile_type: "cesium_ion", cesium_ion_asset_id: 2 }],
    };

    const result = applyFallbacks(sceneProperty);

    expect(result?.default).toEqual(sceneProperty.default);
  });

  it("should not fallback non-cesium_ion tiles", () => {
    const sceneProperty: SceneProperty = {
      tiles: [
        { id: "tile-1", tile_type: "open_street_map" },
        { id: "tile-2", tile_type: "url" },
      ],
    };

    const result = applyFallbacks(sceneProperty);

    expect(result?.tiles).toEqual(sceneProperty.tiles);
  });

  it("should handle mixed tile types", () => {
    const sceneProperty: SceneProperty = {
      tiles: [
        { id: "tile-1", tile_type: "cesium_ion", cesium_ion_asset_id: 2 },
        { id: "tile-2", tile_type: "open_street_map" },
        { id: "tile-3", tile_type: "cesium_ion", cesium_ion_asset_id: 3812 },
      ],
    };

    const result = applyFallbacks(sceneProperty);

    expect(result?.tiles).toEqual([
      { id: "tile-1", tile_type: "google_satellite", cesium_ion_asset_id: 2 },
      { id: "tile-2", tile_type: "open_street_map" },
      { id: "tile-3", tile_type: "nasa_black_marble", cesium_ion_asset_id: 3812 },
    ]);
  });

  it("should treat empty string ion token as no token", () => {
    const sceneProperty: SceneProperty = {
      default: {
        ion: "",
      },
      tiles: [{ id: "tile-1", tile_type: "cesium_ion", cesium_ion_asset_id: 2 }],
    };

    const result = applyFallbacks(sceneProperty);

    // Should apply fallback because empty string is falsy
    expect(result?.tiles?.[0].tile_type).toBe("google_satellite");
  });
});

describe("applyLayerFallbacks", () => {
  it("should return undefined when rootLayer is undefined", () => {
    const result = applyLayerFallbacks(undefined, false);
    expect(result).toBeUndefined();
  });

  it("should NOT apply fallbacks when Cesium Ion token is provided", () => {
    const rootLayer: Layer = {
      id: "root",
      extensionId: "tileset",
      property: {
        default: {
          sourceType: "osm",
        },
      },
    };

    const result = applyLayerFallbacks(rootLayer, true);

    // Should return the same layer without fallbacks
    expect(result?.property?.default?.sourceType).toBe("osm");
  });

  it("should fallback osm sourceType to reearth-buildings when no token", () => {
    const rootLayer: Layer = {
      id: "layer-1",
      extensionId: "tileset",
      property: {
        default: {
          sourceType: "osm",
        },
      },
    };

    const result = applyLayerFallbacks(rootLayer, false);

    expect(result?.property?.default?.sourceType).toBe("reearth-buildings");
  });

  it("should not fallback non-tileset layers", () => {
    const rootLayer: Layer = {
      id: "layer-1",
      extensionId: "marker",
      property: {
        default: {
          sourceType: "osm",
        },
      },
    };

    const result = applyLayerFallbacks(rootLayer, false);

    expect(result?.property?.default?.sourceType).toBe("osm");
  });

  it("should not fallback tileset layers with non-osm sourceType", () => {
    const rootLayer: Layer = {
      id: "layer-1",
      extensionId: "tileset",
      property: {
        default: {
          sourceType: "google-photorealistic",
        },
      },
    };

    const result = applyLayerFallbacks(rootLayer, false);

    expect(result?.property?.default?.sourceType).toBe("google-photorealistic");
  });

  it("should recursively apply fallbacks to child layers", () => {
    const rootLayer: Layer = {
      id: "root",
      extensionId: "group",
      children: [
        {
          id: "child-1",
          extensionId: "tileset",
          property: {
            default: {
              sourceType: "osm",
            },
          },
        },
        {
          id: "child-2",
          extensionId: "tileset",
          property: {
            default: {
              sourceType: "osm",
            },
          },
        },
      ],
    };

    const result = applyLayerFallbacks(rootLayer, false);

    expect(result?.children?.[0].property?.default?.sourceType).toBe("reearth-buildings");
    expect(result?.children?.[1].property?.default?.sourceType).toBe("reearth-buildings");
  });

  it("should handle mixed layer types in children", () => {
    const rootLayer: Layer = {
      id: "root",
      extensionId: "group",
      children: [
        {
          id: "child-1",
          extensionId: "tileset",
          property: {
            default: {
              sourceType: "osm",
            },
          },
        },
        {
          id: "child-2",
          extensionId: "marker",
          property: {
            default: {
              sourceType: "osm",
            },
          },
        },
        {
          id: "child-3",
          extensionId: "tileset",
          property: {
            default: {
              sourceType: "google-photorealistic",
            },
          },
        },
      ],
    };

    const result = applyLayerFallbacks(rootLayer, false);

    expect(result?.children?.[0].property?.default?.sourceType).toBe("reearth-buildings");
    expect(result?.children?.[1].property?.default?.sourceType).toBe("osm"); // Not a tileset
    expect(result?.children?.[2].property?.default?.sourceType).toBe("google-photorealistic"); // Not osm
  });

  it("should not mutate original layer", () => {
    const rootLayer: Layer = {
      id: "layer-1",
      extensionId: "tileset",
      property: {
        default: {
          sourceType: "osm",
        },
      },
    };

    const result = applyLayerFallbacks(rootLayer, false);

    // Original should not be modified
    expect(rootLayer.property?.default?.sourceType).toBe("osm");
    // Result should be modified
    expect(result?.property?.default?.sourceType).toBe("reearth-buildings");
  });

  it("should preserve other layer properties", () => {
    const rootLayer: Layer = {
      id: "layer-1",
      extensionId: "tileset",
      isVisible: true,
      property: {
        default: {
          sourceType: "osm",
          tileset: "custom-tileset",
          shadows: "enabled",
        },
      },
    };

    const result = applyLayerFallbacks(rootLayer, false);

    expect(result?.id).toBe("layer-1");
    expect(result?.extensionId).toBe("tileset");
    expect(result?.isVisible).toBe(true);
    expect(result?.property?.default?.sourceType).toBe("reearth-buildings");
    expect(result?.property?.default?.tileset).toBe("custom-tileset");
    expect(result?.property?.default?.shadows).toBe("enabled");
  });

  it("should handle deeply nested layer trees", () => {
    const rootLayer: Layer = {
      id: "root",
      extensionId: "group",
      children: [
        {
          id: "level1",
          extensionId: "group",
          children: [
            {
              id: "level2",
              extensionId: "tileset",
              property: {
                default: {
                  sourceType: "osm",
                },
              },
            },
          ],
        },
      ],
    };

    const result = applyLayerFallbacks(rootLayer, false);

    expect(result?.children?.[0].children?.[0].property?.default?.sourceType).toBe(
      "reearth-buildings",
    );
  });
});
