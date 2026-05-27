export type CustomProviders = {
  imagery?: {
    providers?: {
      id: string;
      url: string;
      credit?: string;
      maximumLevel?: number;
      minimumLevel?: number;
    }[];
  };
  terrain?: {
    providers?: {
      id: string;
      url: string;
      requestVertexNormals?: boolean;
      requestWaterMask?: boolean;
      credit?: string;
    }[];
  };
  layers?: {
    providers?: {
      id: string;
      url: string;
      options?: Record<string, unknown>;
    }[];
  };
};
