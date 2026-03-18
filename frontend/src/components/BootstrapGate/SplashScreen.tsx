// components/SplashScreen.tsx
import styled from "@emotion/styled";
import { keyframes } from "@emotion/react";
import { glassPerfect } from "./styles/glass";

const fadeIn = keyframes`
  from { opacity: 0; transform: scale(0.96); }
  to { opacity: 1; transform: scale(1); }
`;

const shimmer = keyframes`
  0% { opacity: 0.6; }
  50% { opacity: 1; }
  100% { opacity: 0.6; }
`;

const Container = styled.div`
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: radial-gradient(circle at top, #0f172a, #020617);
`;

const GlassCard = styled.div`
  ${glassPerfect};
  padding: 60px 80px;
  text-align: center;
  animation: ${fadeIn} 0.6s ease-out;
`;

const Logo = styled.h1`
  font-size: 32px;
  letter-spacing: 2px;
  font-weight: 600;
  color: #fff;
`;

const Loader = styled.div`
  margin-top: 20px;
  height: 6px;
  width: 120px;
  border-radius: 999px;
  background: linear-gradient(90deg, #d4af37, #ffb84d);
  animation: ${shimmer} 1.4s infinite ease-in-out;
`;

export default function SplashScreen() {
  return (
    <Container>
      <GlassCard>
        <Logo>ANKHORA</Logo>
        <Loader />
      </GlassCard>
    </Container>
  );
}