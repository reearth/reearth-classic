import React from "react";

// Components
import MoleculeHeader from "@reearth/classic/components/molecules/EarthEditor/Header";
import PublicationModal from "@reearth/classic/components/molecules/EarthEditor/PublicationModal";

import useHooks from "./hooks";

// TODO: ErrorBoudaryでエラーハンドリング

export type Props = {
  className?: string;
};

const Header: React.FC<Props> = ({ className }) => {
  const {
    workspaces = [],
    workspaceId,
    publicationModalVisible,
    searchIndex,
    publishing,
    workspaceModalVisible,
    projectId,
    projectAlias,
    projectStatus,
    publicationModalLoading,
    user,
    currentWorkspace,
    currentProject,
    validAlias,
    validatingAlias,
    url,
    devPluginExtensions,
    handlePublicationModalOpen,
    handlePublicationModalClose,
    handleWorkspaceModalOpen,
    handleWorkspaceModalClose,
    handleSearchIndexChange,
    handleTeamChange,
    handleTeamCreate,
    handleProjectPublish,
    handleProjectAliasCheck,
    handleCopyToClipBoard,
    handlePreviewOpen,
    handleLogout,
    handleDevPluginExtensionsReload,
    handleDevPluginsInstall,
  } = useHooks();

  return (
    <>
      <MoleculeHeader
        className={className}
        currentProjectStatus={projectStatus}
        currentProject={currentProject}
        user={user}
        workspaces={workspaces}
        workspaceId={workspaceId}
        currentWorkspace={currentWorkspace}
        modalShown={workspaceModalVisible}
        devPluginExtensions={devPluginExtensions}
        onPublishmentStatusClick={handlePublicationModalOpen}
        onSignOut={handleLogout}
        onWorkspaceCreate={handleTeamCreate}
        onWorkspaceChange={handleTeamChange}
        onPreviewOpen={handlePreviewOpen}
        openModal={handleWorkspaceModalOpen}
        onModalClose={handleWorkspaceModalClose}
        onDevPluginExtensionsReload={handleDevPluginExtensionsReload}
        onDevPluginsInstall={handleDevPluginsInstall}
      />
      <PublicationModal
        className={className}
        loading={publicationModalLoading}
        isVisible={!!publicationModalVisible}
        searchIndex={searchIndex}
        publishing={publishing}
        projectId={projectId}
        projectAlias={projectAlias}
        publicationStatus={projectStatus}
        validAlias={validAlias}
        validatingAlias={validatingAlias}
        url={url}
        onClose={handlePublicationModalClose}
        onSearchIndexChange={handleSearchIndexChange}
        onPublish={handleProjectPublish}
        onCopyToClipBoard={handleCopyToClipBoard}
        onAliasValidate={handleProjectAliasCheck}
      />
    </>
  );
};

export default Header;
