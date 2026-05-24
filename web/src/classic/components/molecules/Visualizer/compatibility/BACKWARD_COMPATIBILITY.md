# Backward Compatibility System

## Overview

The backward compatibility system ensures that old tile and terrain provider types are automatically migrated to new provider names. This allows existing projects to continue working seamlessly when provider names change.

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

**Note:** See [FALLBACKS.md](./FALLBACKS.md) for information about the fallback system that runs after backward compatibility.

## Implementation

The backward compatibility logic is implemented in three files:

### 1. `backwardCompatibility.ts`
Contains the core migration logic:
- `applyBackwardCompatibility()` - Main entry point that applies all migrations
- `migrateTileType()` - Migrates tile types
- `migrateTerrainType()` - Migrates terrain types

### 2. `backwardCompatibility.test.ts`
Comprehensive unit tests covering:
- All tile type migrations
- All terrain type migrations
- Edge cases (undefined values, missing fields, etc.)
- Immutability (original data is not modified)

### 3. `hooks.ts`
Integration point where backward compatibility is applied:
```typescript
const [overriddenScenePropertyRaw, overrideSceneProperty] = useOverriddenProperty(sceneProperty);

const overriddenSceneProperty = useMemo(
  () => applyBackwardCompatibility(overriddenScenePropertyRaw),
  [overriddenScenePropertyRaw],
);
```

## Tile Type Migrations

### Legacy → New Provider Mappings

| Old Type          | New Type         | Asset ID | Notes                          |
|-------------------|------------------|----------|--------------------------------|
| missing or `default` | `cesium_ion`  | 2        | Cesium World Imagery           |
| `default_label`   | `cesium_ion`     | 3        | Cesium World Imagery + Labels  |
| `default_road`    | `cesium_ion`     | 4        | Cesium World Imagery + Roads   |
| `black_marble`    | `cesium_ion`     | 3812     | NASA Black Marble              |
| `stamen_toner`    | `carto_light`    | -        | Carto Light basemap            |
| `esri_world_topo` | `open_street_map`| -        | OpenStreetMap                  |

**Note:** Tiles without a `tile_type` field are treated the same as `"default"`.

### Example

**Before migration:**
```json
{
  "tiles": [
    { "id": "tile-1", "tile_type": "default" },
    { "id": "tile-2", "tile_type": "stamen_toner" },
    { "id": "tile-3", "tile_url": "https://..." }
  ]
}
```

**After migration:**
```json
{
  "tiles": [
    { "id": "tile-1", "tile_type": "cesium_ion", "cesium_ion_asset_id": 2 },
    { "id": "tile-2", "tile_type": "carto_light" },
    { "id": "tile-3", "tile_url": "https://...", "tile_type": "cesium_ion", "cesium_ion_asset_id": 2 }
  ]
}
```

Note how `tile-3` without a `tile_type` gets migrated to `cesium_ion` with asset ID 2.

## Terrain Type Migrations

### Rules

1. **Missing terrainType:** If `terrain` is enabled but `terrainType` is not set, default to `"cesium"`
2. **ArcGIS migration:** If `terrainType` is `"arcgis"`, change to `"reearth_terrain"`

### Examples

**Missing terrainType:**
```json
// Before
{ "terrain": true }

// After
{ "terrain": true, "terrainType": "cesium" }
```

**ArcGIS migration:**
```json
// Before
{ "terrain": true, "terrainType": "arcgis" }

// After
{ "terrain": true, "terrainType": "reearth_terrain" }
```

## Testing

Run the unit tests:
```bash
yarn test backwardCompatibility.test.ts
```

All 20 tests should pass, covering:
- 9 tests for tile type migrations
- 5 tests for terrain type migrations
- 6 tests for the main `applyBackwardCompatibility()` function

## Adding New Migrations

To add a new backward compatibility rule:

1. **Update the migration function** in `backwardCompatibility.ts`:
   ```typescript
   case "old_type":
     return {
       ...tile,
       tile_type: "new_type",
       // Add any additional fields
     };
   ```

2. **Add test cases** in `backwardCompatibility.test.ts`:
   ```typescript
   it('should migrate "old_type" to "new_type"', () => {
     const tile = { id: "test", tile_type: "old_type" };
     const result = migrateTileType(tile);
     expect(result.tile_type).toBe("new_type");
   });
   ```

3. **Update this documentation** with the new migration rule

## Notes

- **Immutability:** The original `sceneProperty` is never modified. All transformations return new objects.
- **Performance:** Uses `useMemo` to avoid re-computing migrations on every render.
- **Deep cloning:** Uses `lodash-es/cloneDeep` to ensure nested objects are properly cloned.
- **Asset ID preservation:** If a tile already has a `cesium_ion_asset_id`, it won't be overwritten.

## Database Migration

The backend also has corresponding database migrations:
- `260525000715_migrate_tile_types_to_new_providers.go` - Migrates tile types in MongoDB
- `260525001345_update_terrain_types.go` - Migrates terrain types in MongoDB

These migrations run automatically on server startup to update existing data in the database.
