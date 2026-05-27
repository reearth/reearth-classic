/**
 * Scene Property Compatibility & Fallback System
 *
 * This module provides backward compatibility and fallback transformations
 * for scene properties (tiles and terrain).
 *
 * @see README.md for complete documentation
 */

export {
  applyBackwardCompatibility,
  migrateTileType,
  migrateTerrainType,
} from "./backwardCompatibility";
export { applyFallbacks, fallbackTileType, fallbackTerrainType } from "./fallbacks";
