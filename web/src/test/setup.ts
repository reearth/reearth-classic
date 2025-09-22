/* eslint-disable @typescript-eslint/no-namespace */
import { type EmotionMatchers, matchers as emotionMatchers } from "@emotion/jest";
import * as domMatchers from "@testing-library/jest-dom/matchers";
import { cleanup } from "@testing-library/react";
import { afterEach, expect, vi } from "vitest";

// Vitest on GitHub Actions requires TransformStream to run tests with Cesium
import "web-streams-polyfill/es2018";

declare global {
  namespace Vi {
    interface JestAssertion<T = any> extends jest.Matchers<void, T>, EmotionMatchers {
      toHaveStyleRule: EmotionMatchers["toHaveStyleRule"];
    }
  }
}

expect.extend(domMatchers);
expect.extend(emotionMatchers as any);

Object.defineProperty(window, "matchMedia", {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(), // deprecated
    removeListener: vi.fn(), // deprecated
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
});

Object.defineProperty(window, "requestIdleCallback", {
  writable: true,
  value: vi.fn(),
});

// Mock Cesium for tests
vi.mock("cesium", async () => {
  const cesiumMock = {
    // Basic types
    Color: vi.fn().mockImplementation(() => ({
      red: 1,
      green: 1,
      blue: 1,
      alpha: 1,
      equals: vi.fn(() => true),
      clone: vi.fn(),
    })),

    // Enums
    ArcType: { NONE: 0, GEODESIC: 1, RHUMB: 2 },
    ScreenSpaceEventType: {
      LEFT_CLICK: 0,
      LEFT_DOUBLE_CLICK: 1,
      LEFT_DOWN: 2,
      LEFT_UP: 3,
      MIDDLE_CLICK: 4,
      MIDDLE_DOWN: 5,
      MIDDLE_UP: 6,
      RIGHT_CLICK: 7,
      RIGHT_DOWN: 8,
      RIGHT_UP: 9,
      WHEEL: 10,
      MOUSE_MOVE: 11,
    },
    ClockRange: {
      UNBOUNDED: 0,
      CLAMPED: 1,
      LOOP_STOP: 2,
    },
    ClockStep: {
      TICK_DEPENDENT: 0,
      SYSTEM_CLOCK_MULTIPLIER: 1,
      SYSTEM_CLOCK: 2,
    },
    SceneMode: {
      MORPHING: 0,
      COLUMBUS_VIEW: 1,
      SCENE2D: 2,
      SCENE3D: 3,
    },

    // Mock other commonly used Cesium classes with static methods
    Cartesian2: vi.fn().mockImplementation((x = 0, y = 0) => ({
      x,
      y,
    })),

    Cartesian3: Object.assign(
      vi.fn().mockImplementation((x = 0, y = 0, z = 0) => ({
        x,
        y,
        z,
        clone: vi.fn().mockImplementation(result => {
          result = result || { x: 0, y: 0, z: 0 };
          result.x = x;
          result.y = y;
          result.z = z;
          result.clone = vi.fn().mockImplementation(result => {
            result = result || { x: 0, y: 0, z: 0 };
            result.x = x;
            result.y = y;
            result.z = z;
            return result;
          });
          return result;
        }),
      })),
      {
        // Static properties with clone method
        UNIT_X: {
          x: 1,
          y: 0,
          z: 0,
          clone: vi.fn().mockImplementation(result => {
            result = result || { x: 0, y: 0, z: 0 };
            result.x = 1;
            result.y = 0;
            result.z = 0;
            result.clone = vi.fn().mockImplementation(result2 => {
              result2 = result2 || { x: 0, y: 0, z: 0 };
              result2.x = 1;
              result2.y = 0;
              result2.z = 0;
              return result2;
            });
            return result;
          }),
        },
        UNIT_Y: {
          x: 0,
          y: 1,
          z: 0,
          clone: vi.fn().mockImplementation(result => {
            result = result || { x: 0, y: 0, z: 0 };
            result.x = 0;
            result.y = 1;
            result.z = 0;
            result.clone = vi.fn().mockImplementation(result2 => {
              result2 = result2 || { x: 0, y: 0, z: 0 };
              result2.x = 0;
              result2.y = 1;
              result2.z = 0;
              return result2;
            });
            return result;
          }),
        },
        UNIT_Z: {
          x: 0,
          y: 0,
          z: 1,
          clone: vi.fn().mockImplementation(result => {
            result = result || { x: 0, y: 0, z: 0 };
            result.x = 0;
            result.y = 0;
            result.z = 1;
            result.clone = vi.fn().mockImplementation(result2 => {
              result2 = result2 || { x: 0, y: 0, z: 0 };
              result2.x = 0;
              result2.y = 0;
              result2.z = 1;
              return result2;
            });
            return result;
          }),
        },

        // Static methods
        multiplyByScalar: vi.fn().mockImplementation((cartesian, scalar, result) => {
          result = result || { x: 0, y: 0, z: 0 };
          result.x = cartesian.x * scalar;
          result.y = cartesian.y * scalar;
          result.z = cartesian.z * scalar;
          return result;
        }),
        add: vi.fn().mockImplementation(() => ({ x: 0, y: 0, z: 0 })),
        subtract: vi.fn().mockImplementation((left, right, result) => {
          result = result || { x: 0, y: 0, z: 0 };
          result.x = (left?.x ?? 0) - (right?.x ?? 0);
          result.y = (left?.y ?? 0) - (right?.y ?? 0);
          result.z = (left?.z ?? 0) - (right?.z ?? 0);
          result.clone = vi.fn().mockImplementation(result2 => {
            result2 = result2 || { x: 0, y: 0, z: 0 };
            result2.x = result.x;
            result2.y = result.y;
            result2.z = result.z;
            return result2;
          });
          return result;
        }),
        normalize: vi.fn().mockImplementation(() => ({ x: 0, y: 0, z: 1 })),
        clone: vi.fn().mockImplementation((cartesian, result) => {
          if (!cartesian) return { x: 0, y: 0, z: 0 };
          result = result || { x: 0, y: 0, z: 0 };
          result.x = cartesian.x ?? 0;
          result.y = cartesian.y ?? 0;
          result.z = cartesian.z ?? 0;
          return result;
        }),
        fromDegrees: vi.fn().mockImplementation((longitude, latitude, height = 0) => ({
          x: longitude,
          y: latitude,
          z: height,
          _longitude: longitude, // Store original for conversion back
          _latitude: latitude,
        })),
        dot: vi.fn().mockImplementation((left, right) => {
          return left.x * right.x + left.y * right.y + left.z * right.z;
        }),
        angleBetween: vi.fn().mockImplementation(() => {
          // Return a mock angle in radians
          return Math.PI / 4; // 45 degrees
        }),
      },
    ),

    Matrix3: Object.assign(
      vi.fn().mockImplementation(() => ({})),
      {
        // Static methods
        fromQuaternion: vi.fn().mockImplementation((quaternion, result) => {
          result = result || {};
          return result;
        }),
        multiplyByVector: vi.fn().mockImplementation((matrix, vector, result) => {
          result = result || { x: 0, y: 0, z: 0 };
          // Simple passthrough for mocking
          result.x = vector.x ?? 0;
          result.y = vector.y ?? 0;
          result.z = vector.z ?? 0;
          return result;
        }),
      },
    ),

    Matrix4: Object.assign(
      vi.fn().mockImplementation(() => ({})),
      {
        // Static methods
        clone: vi.fn().mockImplementation((matrix, result) => {
          result = result || {};
          return { ...matrix, ...result };
        }),
      },
    ),

    Quaternion: Object.assign(
      vi.fn().mockImplementation(() => ({})),
      {
        // Static methods
        fromAxisAngle: vi.fn().mockImplementation((axis, angle, result) => {
          result = result || {};
          return result;
        }),
      },
    ),

    JulianDate: Object.assign(
      vi.fn().mockImplementation(() => ({
        dayNumber: 0,
        secondsOfDay: 0,
      })),
      {
        // Static methods
        fromIso8601: vi.fn().mockImplementation(iso8601String => ({
          dayNumber: 0,
          secondsOfDay: 0,
          _iso8601String: iso8601String, // store for toDate conversion
        })),
        fromDate: vi.fn().mockImplementation(date => ({
          dayNumber: 0,
          secondsOfDay: 0,
          _date: date, // store for toDate conversion
        })),
        toDate: vi.fn().mockImplementation(julianDate => {
          if (julianDate?._iso8601String) {
            return new Date(julianDate._iso8601String);
          }
          if (julianDate?._date) {
            return julianDate._date;
          }
          return new Date();
        }),
      },
    ),

    // Terrain and geometry classes
    EllipsoidTerrainProvider: vi.fn().mockImplementation(() => ({})),
    Plane: vi.fn().mockImplementation(() => ({
      normal: { x: 0, y: 0, z: 1 },
      distance: 0,
    })),

    // Globe and Ellipsoid classes
    Globe: vi.fn().mockImplementation(ellipsoid => ({
      ellipsoid: ellipsoid || {
        radii: { x: 6378137.0, y: 6378137.0, z: 6356752.314245179 },
        geodeticSurfaceNormal: vi.fn().mockImplementation(() => ({ x: 0, y: 0, z: 1 })),
        cartesianToCartographic: vi.fn().mockImplementation(cartesian => ({
          longitude: ((cartesian._longitude ?? cartesian.x ?? 0) * Math.PI) / 180, // Return radians for CesiumMath.toDegrees
          latitude: ((cartesian._latitude ?? cartesian.y ?? 0) * Math.PI) / 180,
          height: cartesian.z || 0,
        })),
      },
    })),
    Ellipsoid: Object.assign(
      vi.fn().mockImplementation(() => ({
        radii: { x: 6378137.0, y: 6378137.0, z: 6356752.314245179 },
        geodeticSurfaceNormal: vi.fn().mockImplementation(() => ({ x: 0, y: 0, z: 1 })),
        cartesianToCartographic: vi.fn().mockImplementation(cartesian => ({
          longitude: ((cartesian._longitude ?? cartesian.x ?? 0) * Math.PI) / 180, // Return radians for CesiumMath.toDegrees
          latitude: ((cartesian._latitude ?? cartesian.y ?? 0) * Math.PI) / 180,
          height: cartesian.z || 0,
        })),
      })),
      {
        WGS84: {
          radii: { x: 6378137.0, y: 6378137.0, z: 6356752.314245179 },
          geodeticSurfaceNormal: vi.fn().mockImplementation(() => ({ x: 0, y: 0, z: 1 })),
          cartesianToCartographic: vi.fn().mockImplementation(cartesian => ({
            longitude: cartesian.x || 0,
            latitude: cartesian.y || 0,
            height: cartesian.z || 0,
          })),
        },
        UNIT_SPHERE: {
          radii: { x: 1, y: 1, z: 1 },
          geodeticSurfaceNormal: vi.fn().mockImplementation(() => ({ x: 0, y: 0, z: 1 })),
          cartesianToCartographic: vi.fn().mockImplementation(cartesian => ({
            longitude: cartesian.x || 0,
            latitude: cartesian.y || 0,
            height: cartesian.z || 0,
          })),
        },
      },
    ),

    // Entity classes
    Entity: vi.fn().mockImplementation(() => ({
      id: "mock-entity",
      name: "Mock Entity",
      position: null,
      orientation: null,
      model: null,
      point: null,
      polyline: null,
      polygon: null,
      properties: null,
    })),

    // Property classes
    PropertyBag: vi.fn().mockImplementation((properties = {}) => {
      const bag = { ...properties };
      return {
        ...bag,
        removeProperty: vi.fn().mockImplementation(name => {
          delete bag[name];
        }),
        addProperty: vi.fn().mockImplementation((name, value) => {
          bag[name] = value;
        }),
        hasProperty: vi.fn().mockImplementation(name => name in bag),
        getValue: vi.fn().mockImplementation(() => bag),
      };
    }),

    // Material classes
    PolylineDashMaterialProperty: vi.fn().mockImplementation(() => ({
      color: null,
      dashLength: 16,
      dashPattern: 255,
    })),

    // Math utilities - export as both Math and for direct access
    Math: {
      toRadians: vi.fn().mockImplementation(degrees => (degrees * Math.PI) / 180),
      toDegrees: vi.fn().mockImplementation(radians => (radians * 180) / Math.PI),
      PI: Math.PI,
      PI_OVER_TWO: Math.PI / 2,
      TWO_PI: Math.PI * 2,
    },

    // Transforms utilities
    Transforms: {
      eastNorthUpToFixedFrame: vi.fn().mockImplementation((origin, ellipsoid, result) => {
        result = result || {};
        return { ...result, origin, ellipsoid };
      }),
    },

    // Default export
    default: {},
  };

  return cesiumMock;
});

// Mock @cesium/engine
vi.mock("@cesium/engine", () => ({
  default: {},
}));

// Mock resium components
vi.mock("resium", () => ({
  Viewer: vi.fn(({ children }) => children),
  Scene: vi.fn(({ children }) => children),
  Globe: vi.fn(() => null),
  Camera: vi.fn(() => null),
  Clock: vi.fn(() => null),
  Entity: vi.fn(() => null),
  EntityDescription: vi.fn(() => null),
  Billboard: vi.fn(() => null),
  Label: vi.fn(() => null),
  Point: vi.fn(() => null),
  Polyline: vi.fn(() => null),
  Polygon: vi.fn(() => null),
  ImageryLayer: vi.fn(() => null),
  UrlTemplateImageryProvider: vi.fn(),
  OpenStreetMapImageryProvider: vi.fn(),
  KmlDataSource: vi.fn(() => ({
    name: "Mock KML DataSource",
    show: true,
    entities: { values: [] },
    load: vi.fn().mockResolvedValue({}),
  })),
  CzmlDataSource: vi.fn(() => ({
    name: "Mock CZML DataSource",
    show: true,
    entities: { values: [] },
    load: vi.fn().mockResolvedValue({}),
  })),
  GeoJsonDataSource: vi.fn(() => ({
    name: "Mock GeoJSON DataSource",
    show: true,
    entities: { values: [] },
    load: vi.fn().mockResolvedValue({}),
  })),
  default: {},
}));

// Mock cesium-dnd
vi.mock("cesium-dnd", () => ({
  default: vi.fn().mockImplementation(() => ({
    enable: vi.fn(),
    disable: vi.fn(),
    destroy: vi.fn(),
  })),
  Context: vi.fn().mockImplementation(() => ({})),
}));

afterEach(cleanup);
