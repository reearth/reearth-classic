import { useContext } from "react";

import { PluginContext } from "./context";

export default (): boolean => {
  const ctx = useContext(PluginContext);
  return !!ctx?.useExperimentalSandbox;
};
