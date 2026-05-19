import { Meta, StoryObj } from "@storybook/react-vite";

import { V, location } from "../storybook";

import Marker from ".";

const meta: Meta<typeof Marker> = {
  title: "classic/molecules/Visualizer/Engine/Cesium/Marker",
  component: Marker,
};
export default meta;
type Story = StoryObj<typeof Marker>;

export const Point: Story = {
  render: args => (
    <V>
      <Marker {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      isVisible: true,
      property: {
        default: {
          location,
          height: location.height,
          style: "point",
          pointColor: "blue",
          pointSize: 50,
        },
      },
    },
  },
};

export const PointWithLabelAndExcluded: Story = {
  render: args => (
    <V>
      <Marker {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      isVisible: true,
      property: {
        default: {
          location,
          height: location.height,
          style: "point",
          pointColor: "blue",
          pointSize: 50,
          extrude: true,
          label: true,
          labelText: "label",
        },
      },
    },
  },
};

export const PointWithRightLabel: Story = {
  render: args => (
    <V>
      <Marker {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      isVisible: true,
      property: {
        default: {
          location,
          height: location.height,
          style: "point",
          label: true,
          labelText: "label",
          labelPosition: "left",
          labelTypography: {
            fontSize: 15,
            color: "red",
            bold: true,
            italic: true,
            fontFamily: "serif",
          },
        },
      },
    },
  },
};

export const PointWithTopLabel: Story = {
  render: args => (
    <V>
      <Marker {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      isVisible: true,
      property: {
        default: {
          location,
          height: location.height,
          style: "point",
          label: true,
          labelText: "label",
          labelPosition: "top",
          labelTypography: {
            fontFamily: "serif",
          },
        },
      },
    },
  },
};

export const PointWithBottomLabel: Story = {
  render: args => (
    <V>
      <Marker {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      isVisible: true,
      property: {
        default: {
          location,
          height: location.height,
          style: "point",
          label: true,
          labelText: "label",
          labelPosition: "bottom",
        },
      },
    },
  },
};

export const Image: Story = {
  render: args => (
    <V>
      <Marker {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      isVisible: true,
      property: {
        default: {
          location,
          height: location.height,
          style: "image",
          image: `/sample.svg`,
        },
      },
    },
  },
};

export const ImageWithShadow: Story = {
  render: args => (
    <V>
      <Marker {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      isVisible: true,
      property: {
        default: {
          location,
          height: location.height,
          style: "image",
          image: `/sample.png`,
          imageShadow: true,
        },
      },
    },
  },
};

export const ImageWithCropAndShadow: Story = {
  render: args => (
    <V>
      <Marker {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      isVisible: true,
      property: {
        default: {
          location,
          height: location.height,
          style: "image",
          image: `/sample.png`,
          imageCrop: "circle",
          imageShadow: true,
          extrude: true,
        },
      },
    },
  },
};

export const ImageWithColor: Story = {
  render: args => (
    <V>
      <Marker {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      isVisible: true,
      property: {
        default: {
          location,
          height: location.height,
          style: "image",
          image: `/sample.png`,
          imageCrop: "circle",
          imageShadow: true,
          extrude: true,
          imageColor: "red",
        },
      },
    },
  },
};

export const ImageWithRightLabel: Story = {
  render: args => (
    <V>
      <Marker {...args} />
    </V>
  ),
  args: {
    layer: {
      id: "",
      isVisible: true,
      property: {
        default: {
          location,
          height: location.height,
          style: "image",
          image: `/sample.png`,
          label: true,
          labelText: "label",
          labelPosition: "right",
          labelTypography: {
            fontSize: 15,
            color: "red",
            bold: true,
            italic: true,
            fontFamily: "serif",
          },
        },
      },
    },
  },
};
