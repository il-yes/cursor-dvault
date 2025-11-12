import { DashboardLayout } from "@/components/DashboardLayout";
import { VaultThreeColumnLayout } from "@/components/VaultThreeColumnLayout";
import { useParams } from "react-router-dom";

const Vault = () => {
  const { filter, folderId } = useParams<{ filter?: string; folderId?: string }>();
  const activeFilter = folderId ? `folder:${folderId}` : (filter || "all");

  return (
    <DashboardLayout>
      <VaultThreeColumnLayout filter={activeFilter} />
    </DashboardLayout>
  );
};

export default Vault;
