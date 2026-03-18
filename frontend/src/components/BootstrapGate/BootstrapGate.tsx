// components/BootstrapGate.tsx
import { useEffect, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import SplashScreen from "./SplashScreen";


import * as AppAPI from "../../../wailsjs/go/main/App";
import OnboardingWizardBeta from "../onBoardingWizardBeta";
import EmailLookup from "@/pages/EmailLookup";
import { useNavigate } from "react-router-dom";

export default function BootstrapGate() {
  const [isLoading, setIsLoading] = useState(true);

//   useEffect(() => {
//     // Must be inside a useEffect, not top-level
//     init();
//   }, [init]);

  return <SplashScreen />;
}

