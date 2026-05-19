/* eslint-disable @typescript-eslint/no-unused-vars */
import { Meta, StoryObj } from "@storybook/react-vite";
import React, { useState } from "react";

import TagPane, { TagGroup } from ".";

const meta: Meta<typeof TagPane> = {
  title: "classic/molecules/EarthEditor/TagPane",
  component: TagPane,
};
export default meta;
type Story = StoryObj<typeof TagPane>;

export const Default: Story = {
  render: () => {
    const [allTagGroups, setAllTagGroups] = useState<TagGroup[]>([
      {
        id: "default",
        label: "Default",
        tags: [
          { id: "hoge", label: "hoge" },
          { id: "fuga", label: "fuga" },
          { id: "foo", label: "foo" },
          { id: "wow", label: "wow" },
        ],
      },
      {
        id: "year",
        label: "Year",
        tags: [
          { id: "1995", label: "1995" },
          { id: "2000", label: "2000" },
          { id: "2005", label: "2005" },
        ],
      },
    ]);
    const [tagGroups, setTagGroups] = useState<TagGroup[]>([
      {
        id: "default",
        label: "Default",
        tags: [
          { id: "hoge", label: "hoge" },
          { id: "fuga", label: "fuga" },
        ],
      },
      {
        id: "year",
        label: "Year",
        tags: [
          { id: "1995", label: "1995" },
          { id: "2000", label: "2000" },
        ],
      },
    ]);
    return (
      <TagPane
        allTagGroups={allTagGroups}
        // onTagGroupAdd={handleAddTagGroup}
        // onTagAdd={handleAddTag}
      />
    );
  },
};
