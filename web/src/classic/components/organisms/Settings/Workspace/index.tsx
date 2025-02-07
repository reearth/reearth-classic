import React, { useCallback, useEffect, useState } from "react";

import Loading from "@reearth/classic/components/atoms/Loading";
import SettingsHeader from "@reearth/classic/components/molecules/Settings/SettingsHeader";
import DangerSection from "@reearth/classic/components/molecules/Settings/Workspace/DangerSection";
import MembersSection from "@reearth/classic/components/molecules/Settings/Workspace/MembersSection";
import ProfileSection from "@reearth/classic/components/molecules/Settings/Workspace/ProfileSection";
import SettingPage from "@reearth/classic/components/organisms/Settings/SettingPage";
import { Workspace } from "@reearth/services/state";

import useHooks from "./hooks";

type Props = {
  workspaceId: string;
};

const WorkspaceSettings: React.FC<Props> = ({ workspaceId }) => {
  const {
    me,
    currentWorkspace,
    currentProject,
    searchedUser,
    changeSearchedUser,
    deleteWorkspace,
    updateName,
    searchUser,
    addMembersToWorkspace,
    updateMemberOfWorkspace,
    removeMemberFromWorkspace,
    loading,
  } = useHooks({ workspaceId });
  const [owner, setOwner] = useState(false);
  const members = currentWorkspace?.members;

  const checkOwner = useCallback(() => {
    if (members) {
      for (let i = 0; i < members.length; i++) {
        if (members[i].userId === me?.id && members[i].role === "OWNER") {
          return true;
        }
      }
    }
    return false;
  }, [members, me?.id]);

  useEffect(() => {
    const o = checkOwner();
    setOwner(o);
  }, [checkOwner]);

  const WorkspaceContent: React.FC<{ currentWorkspace: Workspace; owner: boolean }> = ({
    currentWorkspace,
    owner,
  }) => (
    <>
      <ProfileSection
        currentWorkspace={currentWorkspace}
        updateWorkspaceName={updateName}
        owner={owner}
      />
      {!currentWorkspace.personal && (
        <MembersSection
          me={me}
          owner={owner}
          members={members}
          searchedUser={searchedUser}
          changeSearchedUser={changeSearchedUser}
          searchUser={searchUser}
          addMembersToWorkspace={addMembersToWorkspace}
          updateMemberOfWorkspace={updateMemberOfWorkspace}
          removeMemberFromWorkspace={removeMemberFromWorkspace}
        />
      )}
    </>
  );

  return (
    <SettingPage workspaceId={workspaceId} projectId={currentProject?.id}>
      <SettingsHeader currentWorkspace={currentWorkspace} />
      {currentWorkspace && <WorkspaceContent currentWorkspace={currentWorkspace} owner={owner} />}
      {me.myTeam && me.myTeam !== workspaceId && (
        <DangerSection workspace={currentWorkspace} deleteWorkspace={deleteWorkspace} />
      )}
      {loading && <Loading portal overlay />}
    </SettingPage>
  );
};

export default WorkspaceSettings;
