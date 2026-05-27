# Scene Property Compatibility & Fallback System

## Overview

This directory contains a comprehensive system for handling backward compatibility and fallbacks for scene properties, particularly for tile and terrain providers. The system ensures old projects continue working while gracefully handling missing authentication credentials.

## System Architecture

### Data Flow Pipeline

```
┌─────────────────────────────────────────────────────────────────────┐
│ Input: sceneProperty (from props)                                   │
└────────────────────────┬────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────────────────┐
│ Step 1: useOverriddenProperty()                                     │
│ - Apply plugin overrides                                            │
└────────────────────────┬────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────────────────┐
│ Step 2: applyBackwardCompatibility()                                │
│ - Migrate old provider types to new types                           │
│ - Example: "default" → "cesium_ion" with asset_id: 2                │
└────────────────────────┬────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────────────────┐
│ Step 3: applyFallbacks()                                            │
│ - Only if no Cesium Ion token present                               │
│ - Fallback to alternative providers                                 │
│ - Example: "cesium_ion" id 2 → "google_satellite"                   │
└────────────────────────┬────────────────────────────────────────────┘
                         ↓
┌─────────────────────────────────────────────────────────────────────┐
│ Output: overriddenSceneProperty (passed to consumers)               │
└─────────────────────────────────────────────────────────────────────┘
```

## Files in This Directory

### Core Implementation

| File | Purpose | Lines | Tests |
|------|---------|-------|-------|
| `backwardCompatibility.ts` | Migrate old types to new types | ~130 | 20 tests |
| `fallbacks.ts` | Fallback when no Cesium Ion token | ~120 | 22 tests |
| `index.ts` | Public API exports | ~10 | N/A |

### Tests

| File | Purpose | Tests | Status |
|------|---------|-------|--------|
| `backwardCompatibility.test.ts` | Unit tests for migrations | 20 | ✅ All pass |
| `fallbacks.test.ts` | Unit tests for fallbacks | 22 | ✅ All pass |
| `compatibility.integration.test.ts` | Integration tests | 7 | ✅ All pass |

**Total: 49 tests, all passing**

### Documentation

| File | Purpose |
|------|---------|
| `BACKWARD_COMPATIBILITY.md` | Detailed backward compatibility guide |
| `FALLBACKS.md` | Detailed fallback system guide |
| `README.md` | This file - system overview |

## Quick Reference

### Backward Compatibility Rules

**Tiles:**
- `"default"` → `cesium_ion` (asset_id: 2)
- `"default_label"` → `cesium_ion` (asset_id: 3)
- `"default_road"` → `cesium_ion` (asset_id: 4)
- `"black_marble"` → `cesium_ion` (asset_id: 3812)
- `"stamen_toner"` → `open_street_map`
- `"esri_world_topo"` → `open_street_map`

**Terrain:**
- `"arcgis"` → `"reearth_terrain"`

### Fallback Rules (Only When No Cesium Ion Token)

**Tiles:**
- `cesium_ion` asset_id 2 → `google_satellite`
- `cesium_ion` asset_id 3 → `google_satellite`
- `cesium_ion` asset_id 4 → `google_roadmap`
- `cesium_ion` asset_id 3812 → `nasa_black_marble`

**Terrain:**
- `"cesium"` → `"reearth_terrain"`

## Usage

### Import from index

```typescript
import { applyBackwardCompatibility, applyFallbacks } from "./compatibility";
```

### Apply in hooks

```typescript
// Step 1: Plugin overrides
const [overriddenScenePropertyRaw, overrideSceneProperty] = useOverriddenProperty(sceneProperty);

// Step 2: Backward compatibility
const backwardCompatibleSceneProperty = useMemo(
  () => applyBackwardCompatibility(overriddenScenePropertyRaw),
  [overriddenScenePropertyRaw],
);

// Step 3: Fallbacks
const overriddenSceneProperty = useMemo(
  () => applyFallbacks(backwardCompatibleSceneProperty),
  [backwardCompatibleSceneProperty],
);
```

## Testing

### Run All Tests
```bash
yarn test compatibility
```

### Test Coverage
- ✅ Backward compatibility migrations
- ✅ Fallback logic with/without token
- ✅ Integration of both systems
- ✅ Edge cases (undefined, empty, mixed types)
- ✅ Immutability verification
- ✅ Property preservation

## Performance Considerations

- **Immutable**: Uses deep cloning to avoid mutating original data
- **Memoized**: Both transformations use `useMemo` in hooks
- **Efficient**: Only processes when inputs change
- **No re-renders**: Transformations don't trigger unnecessary re-renders

## Future Work

### Short Term
- Monitor fallback usage metrics
- Gradually migrate users to proper Cesium Ion tokens

### Long Term (When Fallbacks Can Be Removed)
1. All users have Cesium Ion tokens configured
2. Alternative providers are fully stable
3. Remove `fallbacks.ts` and related code
4. Keep backward compatibility for historical data

## Database Migrations

Corresponding server-side migrations in MongoDB:

| Migration | File | Purpose |
|-----------|------|---------|
| 260525000715 | `migrate_tile_types_to_new_providers.go` | Migrate tile types in database |
| 260525001345 | `update_terrain_types.go` | Migrate terrain types in database |

These migrations run automatically on server startup.

## Support & Documentation

- **Backward Compatibility**: See [BACKWARD_COMPATIBILITY.md](./BACKWARD_COMPATIBILITY.md)
- **Fallbacks**: See [FALLBACKS.md](./FALLBACKS.md)
- **Questions**: Check existing tests for usage examples

## Status

✅ **Production Ready**
- All 49 tests passing
- Comprehensive documentation
- Integration tested
- Build verified
- Properly organized in subfolder
