# Fallback System

## Overview

The fallback system provides temporary alternatives when Cesium Ion assets cannot be used due to missing authentication tokens. This is a **temporary measure** and will be removed in the future once all users have proper Cesium Ion token configuration.

## Data Flow

```
sceneProperty (from props)
    ↓
useOverriddenProperty() - Apply plugin overrides
    ↓
applyBackwardCompatibility() - Migrate old types to new types
    ↓
applyFallbacks() - Fallback when no Cesium Ion token
    ↓
overriddenSceneProperty (passed to consumers)
```

## When Fallbacks Apply

Fallbacks **ONLY** apply when:
- No Cesium Ion token is provided (`sceneProperty.default?.ion` is empty/undefined)
- The tile or terrain type requires Cesium Ion authentication

If a Cesium Ion token is present, no fallbacks are applied and the original Cesium Ion assets are used.

## Implementation

The fallback logic is implemented in three files:

### 1. `fallbacks.ts`
Contains the core fallback logic:
- `applyFallbacks()` - Main entry point that checks for token and applies fallbacks
- `fallbackTileType()` - Fallbacks for tiles requiring Cesium Ion
- `fallbackTerrainType()` - Fallbacks for terrain requiring Cesium Ion

### 2. `fallbacks.test.ts`
Comprehensive unit tests covering:
- All tile type fallbacks
- Terrain type fallback
- Token presence/absence scenarios
- Edge cases (empty token, mixed tile types, etc.)

### 3. `hooks.ts`
Integration point where fallbacks are applied after backward compatibility:
```typescript
const backwardCompatibleSceneProperty = useMemo(
  () => applyBackwardCompatibility(overriddenScenePropertyRaw),
  [overriddenScenePropertyRaw],
);

const overriddenSceneProperty = useMemo(
  () => applyFallbacks(backwardCompatibleSceneProperty),
  [backwardCompatibleSceneProperty],
);
```

## Tile Fallback Rules

### When No Cesium Ion Token

| Cesium Ion Asset ID | Fallback Provider   | Notes                          |
|---------------------|---------------------|--------------------------------|
| 2                   | `google_satellite`  | Google Satellite Imagery       |
| 3                   | `google_satellite`  | Google Satellite Imagery       |
| 4                   | `google_roadmap`    | Google Road Map                |
| 3812                | `nasa_black_marble` | NASA Black Marble              |

**Note:** Other Cesium Ion asset IDs remain unchanged (no fallback).

### Example

**Scenario: No Cesium Ion token**

Before fallback:
```json
{
  "default": {},
  "tiles": [
    { "id": "tile-1", "tile_type": "cesium_ion", "cesium_ion_asset_id": 2 },
    { "id": "tile-2", "tile_type": "cesium_ion", "cesium_ion_asset_id": 4 }
  ]
}
```

After fallback:
```json
{
  "default": {},
  "tiles": [
    { "id": "tile-1", "tile_type": "google_satellite", "cesium_ion_asset_id": 2 },
    { "id": "tile-2", "tile_type": "google_roadmap", "cesium_ion_asset_id": 4 }
  ]
}
```

**Scenario: With Cesium Ion token**

```json
{
  "default": { "ion": "my-token-here" },
  "tiles": [
    { "id": "tile-1", "tile_type": "cesium_ion", "cesium_ion_asset_id": 2 }
  ]
}
```

No fallback applied - tiles remain as `cesium_ion` with original asset IDs.

## Terrain Fallback Rules

### When No Cesium Ion Token

| Terrain Type | Fallback Type      | Notes                          |
|--------------|--------------------|---------------------------------|
| `cesium`     | `reearth_terrain`  | Re:Earth terrain service       |

**Note:** Other terrain types remain unchanged (no fallback).

### Example

**Scenario: No Cesium Ion token**

Before fallback:
```json
{
  "terrain": {
    "terrain": true,
    "terrainType": "cesium"
  }
}
```

After fallback:
```json
{
  "terrain": {
    "terrain": true,
    "terrainType": "reearth_terrain"
  }
}
```

## Combined Example: Backward Compatibility + Fallbacks

### Input (Old Format)
```json
{
  "default": {},  // No Cesium Ion token
  "tiles": [
    { "id": "tile-1", "tile_type": "default" },
    { "id": "tile-2", "tile_type": "black_marble" }
  ],
  "terrain": {
    "terrain": true
  }
}
```

### After Backward Compatibility
```json
{
  "default": {},
  "tiles": [
    { "id": "tile-1", "tile_type": "cesium_ion", "cesium_ion_asset_id": 2 },
    { "id": "tile-2", "tile_type": "cesium_ion", "cesium_ion_asset_id": 3812 }
  ],
  "terrain": {
    "terrain": true,
    "terrainType": "cesium"
  }
}
```

### After Fallbacks (Final Result)
```json
{
  "default": {},
  "tiles": [
    { "id": "tile-1", "tile_type": "google_satellite", "cesium_ion_asset_id": 2 },
    { "id": "tile-2", "tile_type": "nasa_black_marble", "cesium_ion_asset_id": 3812 }
  ],
  "terrain": {
    "terrain": true,
    "terrainType": "reearth_terrain"
  }
}
```

## Testing

Run the unit tests:
```bash
yarn test fallbacks.test.ts
```

All 22 tests should pass, covering:
- 8 tests for tile type fallbacks
- 4 tests for terrain type fallback
- 10 tests for the main `applyFallbacks()` function

## Future Removal

This fallback system is temporary and will be removed in a future release when:
1. All users have been migrated to proper Cesium Ion token configuration
2. Alternative providers are fully established and tested
3. Backward compatibility migration period has ended

## Notes

- **Token Check**: Fallbacks only apply when `sceneProperty.default?.ion` is falsy (undefined, null, or empty string)
- **Immutability**: The original `sceneProperty` is never modified. All transformations return new objects.
- **Performance**: Uses `useMemo` to avoid re-computing fallbacks on every render.
- **Deep cloning**: Uses `lodash-es/cloneDeep` to ensure nested objects are properly cloned.
- **Selective Fallback**: Only specific Cesium Ion asset IDs are mapped to fallbacks. Unknown asset IDs remain unchanged.
- **No Double Fallback**: If backward compatibility already changed the type to something other than `cesium_ion`, fallbacks won't apply.
