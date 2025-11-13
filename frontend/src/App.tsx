import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Home from "./pages/Home";
import Index from "./pages/Index";
import Vault from "./pages/Vault";
import OfflineVault from "./pages/OfflineVault";
import SignIn from "./pages/SignIn";
import NotFound from "./pages/NotFound";
import ShareEntries from "./pages/ShareEntries";

const queryClient = new QueryClient();

const App = () => (
  <QueryClientProvider client={queryClient}>
    <Toaster />
    <Sonner />
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/dashboard" element={<Index />} />
        <Route path="/dashboard/vault" element={<Vault />} />
        <Route path="/dashboard/vault/:filter" element={<Vault />} />
        <Route path="/dashboard/vault/folder/:folderId" element={<Vault />} />
        <Route path="/dashboard/shared" element={<ShareEntries />} />
        <Route path="/vault/offline" element={<OfflineVault />} />
        <Route path="/auth/signin" element={<SignIn />} />
        {/* ADD ALL CUSTOM ROUTES ABOVE THE CATCH-ALL "*" ROUTE */}
        <Route path="*" element={<NotFound />} />
      </Routes>
    </BrowserRouter>
  </QueryClientProvider>
);

export default App;
