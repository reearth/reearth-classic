/* eslint-disable @typescript-eslint/no-unused-vars */
import { Meta, StoryObj } from "@storybook/react-vite";
import React, { useState } from "react";

import TagGroup from ".";

const meta: Meta<typeof TagGroup> = {
  title: "classic/molecules/EarthEditor/TagGroup/TagGroup",
  component: TagGroup,
};
export default meta;
type Story = StoryObj<typeof TagGroup>;

const DefaultRenderer = () => {
  const [attachedTags, setAttachedTags] = useState([
    { id: "hoge", label: "hoge" },
    { id: "fuga", label: "fuga" },
  ]);
  const [allTags, setAllTags] = useState([
    { id: "hoge", label: "hoge" },
    { id: "fuga", label: "fuga" },
    { id: "foo", label: "foo" },
  ]);
  return (
    <TagGroup
      attachedTags={attachedTags}
      allTags={allTags}
      // onSelect={handleSelect}
      // onTagAdd={handleCreate}
      // onTagRemove={handleDetach}
      title="Default"
      icon="cancel"
    />
  );
};

export const Default: Story = {
  render: () => <DefaultRenderer />,
};
