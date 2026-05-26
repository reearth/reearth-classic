import { describe, expect, it } from "vitest";

import type { SceneProperty } from "../Engine/ref";

import { applyBackwardCompatibility } from "./backwardCompatibility";
import { applyFallbacks } from "./fallbacks";

describe("Backward Compatibility + Fallbacks Integration", () => {
  it("should apply both backward compatibility and fallbacks when no token", () => {
    // Old format with no Cesium Ion token
    const input: SceneProperty = {
      default: {},
      tiles: [
        { id: "tile-1", tile_type: "default" },
        { id: "tile-2", tile_type: "black_marble" },
        { id: "tile-3", tile_type: "stamen_toner" },
      ],
      terrain: {
        terrain: true,
      },
    };

    // Step 1: Apply backward compatibility
    const afterBackwardCompat = applyBackwardCompatibility(input);

    expect(afterBackwardCompat?.tiles).toEqual([
      { id: "tile-1", tile_type: "cesium_ion", cesium_ion_asset_id: 2 },
      { id: "tile-2", tile_type: "cesium_ion", cesium_ion_asset_id: 3812 },
      { id: "tile-3", tile_type: "open_street_map" },
    ]);
    expect(afterBackwardCompat?.terrain).toEqual({
      terrain: true,
      // terrainType is NOT added when missing
    });

    // Step 2: Apply fallbacks (no token, so fallbacks apply)
    const final = applyFallbacks(afterBackwardCompat);

    expect(final?.tiles).toEqual([
      { id: "tile-1", tile_type: "google_satellite", cesium_ion_asset_id: 2 },
      { id: "tile-2", tile_type: "nasa_black_marble", cesium_ion_asset_id: 3812 },
      { id: "tile-3", tile_type: "open_street_map" }, // No change (not cesium_ion)
    ]);
    expect(final?.terrain).toEqual({
      terrain: true,
      // terrainType remains undefined (no fallback when missing)
    });
  });

  it("should apply backward compatibility but NOT fallbacks when token is present", () => {
    // Old format with Cesium Ion token
    const input: SceneProperty = {
      default: {
        ion: "my-cesium-ion-token",
      },
      tiles: [
        { id: "tile-1", tile_type: "default" },
        { id: "tile-2", tile_type: "default_label" },
      ],
      terrain: {
        terrain: true,
      },
    };

    // Step 1: Apply backward compatibility
    const afterBackwardCompat = applyBackwardCompatibility(input);

    expect(afterBackwardCompat?.tiles).toEqual([
      { id: "tile-1", tile_type: "cesium_ion", cesium_ion_asset_id: 2 },
      { id: "tile-2", tile_type: "cesium_ion", cesium_ion_asset_id: 3 },
    ]);
    expect(afterBackwardCompat?.terrain).toEqual({
      terrain: true,
      // terrainType is NOT added when missing
    });

    // Step 2: Apply fallbacks (has token, so NO fallbacks)
    const final = applyFallbacks(afterBackwardCompat);

    // Should remain as cesium_ion because token is present
    expect(final?.tiles).toEqual([
      { id: "tile-1", tile_type: "cesium_ion", cesium_ion_asset_id: 2 },
      { id: "tile-2", tile_type: "cesium_ion", cesium_ion_asset_id: 3 },
    ]);
    expect(final?.terrain).toEqual({
      terrain: true,
      // terrainType remains undefined (token doesn't affect missing terrainType)
    });
  });

  it("should handle tile without tile_type (no migration)", () => {
    const input: SceneProperty = {
      tiles: [{ id: "tile-1", tile_url: "https://example.com" }],
    };

    // Step 1: Backward compatibility (does NOT add tile_type when missing)
    const afterBackwardCompat = applyBackwardCompatibility(input);

    expect(afterBackwardCompat?.tiles?.[0]).toEqual({
      id: "tile-1",
      tile_url: "https://example.com",
      // tile_type and cesium_ion_asset_id are NOT added
    });

    // Step 2: Apply fallbacks (no tile_type, so no fallback applies)
    const final = applyFallbacks(afterBackwardCompat);

    expect(final?.tiles?.[0]).toEqual({
      id: "tile-1",
      tile_url: "https://example.com",
      // Remains unchanged
    });
  });

  it("should handle esri_world_topo migration without fallback", () => {
    const input: SceneProperty = {
      tiles: [{ id: "tile-1", tile_type: "esri_world_topo" }],
    };

    // Step 1: Backward compatibility
    const afterBackwardCompat = applyBackwardCompatibility(input);

    expect(afterBackwardCompat?.tiles?.[0]).toEqual({
      id: "tile-1",
      tile_type: "open_street_map",
    });

    // Step 2: Fallbacks (no cesium_ion, so no fallback needed)
    const final = applyFallbacks(afterBackwardCompat);

    expect(final?.tiles?.[0]).toEqual({
      id: "tile-1",
      tile_type: "open_street_map",
    });
  });

  it("should handle arcgis (legacy) terrain migration without fallback", () => {
    const input: SceneProperty = {
      terrain: {
        terrain: true,
        terrainType: "arcgis" as any, // Legacy type not in current union
      },
    };

    // Step 1: Backward compatibility
    const afterBackwardCompat = applyBackwardCompatibility(input);

    expect(afterBackwardCompat?.terrain).toEqual({
      terrain: true,
      terrainType: "reearth_terrain",
    });

    // Step 2: Fallbacks (not "cesium", so no fallback applies)
    const final = applyFallbacks(afterBackwardCompat);

    expect(final?.terrain).toEqual({
      terrain: true,
      terrainType: "reearth_terrain",
    });
  });

  it("should handle mixed scenario: some tiles need fallback, some don't", () => {
    const input: SceneProperty = {
      tiles: [
        { id: "tile-1", tile_type: "default" }, // → cesium_ion → google_satellite
        { id: "tile-2", tile_type: "esri_world_topo" }, // → open_street_map (no fallback)
        { id: "tile-3", tile_type: "default_road" }, // → cesium_ion → google_roadmap
      ],
    };

    // Apply both transformations
    const afterBackwardCompat = applyBackwardCompatibility(input);
    const final = applyFallbacks(afterBackwardCompat);

    expect(final?.tiles).toEqual([
      { id: "tile-1", tile_type: "google_satellite", cesium_ion_asset_id: 2 },
      { id: "tile-2", tile_type: "open_street_map" },
      { id: "tile-3", tile_type: "google_roadmap", cesium_ion_asset_id: 4 },
    ]);
  });

  it("should preserve properties through both transformations", () => {
    const input: SceneProperty = {
      default: {
        camera: { lat: 35, lng: 139, height: 1000, heading: 0, pitch: -45, roll: 0, fov: 1 },
        bgcolor: "#1a1a1a",
      },
      tiles: [{ id: "tile-1", tile_type: "default", tile_opacity: 0.8 }],
      atmosphere: {
        enable_sun: true,
        enable_lighting: true,
      },
    };

    const afterBackwardCompat = applyBackwardCompatibility(input);
    const final = applyFallbacks(afterBackwardCompat);

    // Non-tile/terrain properties should be preserved
    expect(final?.default).toEqual(input.default);
    expect(final?.atmosphere).toEqual(input.atmosphere);

    // Tile should be transformed with opacity preserved
    expect(final?.tiles?.[0]).toEqual({
      id: "tile-1",
      tile_type: "google_satellite",
      cesium_ion_asset_id: 2,
      tile_opacity: 0.8,
    });
  });
});
