import { useState } from 'react';


function OnboardingWizard({ onComplete }) {
  const [step, setStep] = useState(1);
  const totalSteps = 3;

  function nextStep() {
    if (step < totalSteps) {
      setStep(step + 1);
    } else {
      onComplete();
    }
  }

  function prevStep() {
    if (step > 1) {
      setStep(step - 1);
    }
  }

  function skipOnboarding() {
    onComplete();
  }

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-gradient-to-br from-white/80 via-white/80 to-zinc-100/90 dark:from-zinc-900/90 dark:via-zinc-900/90 dark:to-black/80 backdrop-blur-3xl z-50">
      <div className="max-w-xl w-full rounded-[2.5rem] border-[1.5px] border-white/40 dark:border-zinc-700/40 bg-gradient-to-br from-white/55 via-white/40 to-zinc-200/25 dark:from-zinc-900/40 dark:to-black/20 shadow-2xl p-12 flex flex-col relative">

        {/* Header */}
        <div className="flex justify-between items-center mb-10 pb-5 border-b border-white/20 dark:border-zinc-700/30">
          <h1 className="text-4xl md:text-5xl font-black tracking-tight bg-gradient-to-r from-[#C9A44A] via-[#D2AF56] to-[#B8934A] bg-clip-text text-transparent drop-shadow-lg">Welcome to Ankhora</h1>
          <button
            onClick={skipOnboarding}
            className="rounded-xl px-5 py-2 text-base font-semibold text-[#C9A44A] bg-white/50 dark:bg-zinc-800/40 backdrop-blur-md hover:bg-white/70 shadow-sm hover:shadow-md transition disabled:opacity-60"
          >
            Skip
          </button>
        </div>

        {/* Content */}
        <div className="flex flex-col items-center text-center mb-8 flex-grow min-h-[320px]">
          {step === 1 && (
            <div className="animate-fadeInUp">
              <div className="w-24 h-24 mb-8 rounded-3xl bg-gradient-to-br from-[#C9A44A]/70 via-[#D2AF56]/60 to-[#B8934A]/80 shadow-2xl flex items-center justify-center text-5xl animate-pulse drop-shadow-lg backdrop-blur-sm">
                🔒
              </div>
              <h2 className="text-2xl md:text-3xl font-bold text-[#C9A44A] mb-4 leading-tight drop-shadow-lg">
                Encryption Happens<br />on Your Device
              </h2>
              <p className="text-lg md:text-xl text-zinc-700 dark:text-zinc-100/90 font-light mb-8 max-w-lg mx-auto leading-relaxed">
                Unlike Dropbox or Google Drive, Ankhora encrypts your files <span className="font-semibold">before they leave your computer</span>. Our servers never see your data in plaintext.
              </p>
              <div className="flex gap-3 justify-center items-center text-[#C9A44A] text-lg max-w-lg mx-auto">
                <span className="font-bold">💻 Your Device</span>
                <span className="font-bold text-[#D2AF56]">→ 🔒 Encrypted →</span>
                <span className="font-bold">☁️ Our Server</span>
              </div>
            </div>
          )}

          {step === 2 && (
            <div className="animate-fadeInUp">
              <div className="w-24 h-24 mb-8 rounded-3xl bg-gradient-to-br from-[#C9A44A]/60 via-[#D2AF56]/60 to-[#B8934A]/70 shadow-2xl flex items-center justify-center text-5xl animate-pulse drop-shadow-lg backdrop-blur-sm">
                🔐
              </div>
              <h2 className="text-2xl md:text-3xl font-bold text-[#C9A44A] mb-4 leading-tight drop-shadow-lg">
                Only You Have the Keys
              </h2>
              <p className="text-lg md:text-xl text-zinc-700 dark:text-zinc-100/90 font-light mb-8 max-w-lg mx-auto leading-relaxed">
                Your encryption keys are derived from your password and stored <span className="font-semibold">only on your device</span>. We never have access to them, which means we can’t decrypt your data—even if we wanted to.
              </p>
              <div className="flex justify-center gap-8 max-w-lg mx-auto flex-wrap">
                <div className="flex flex-col items-center bg-[#C9A44A]/30 backdrop-blur-sm rounded-2xl p-6 shadow-md text-[#C9A44A] font-bold">
                  <span className="text-3xl mb-2">👤</span>
                  <span className="mb-1">You</span>
                  <span className="text-[#D2AF56] font-extrabold text-lg">🔑 Has Key</span>
                </div>
                <div className="flex flex-col items-center bg-[#C9A44A]/10 backdrop-blur-sm rounded-2xl p-6 text-[#C9A44A] text-opacity-70 font-medium">
                  <span className="text-3xl mb-2">🏢</span>
                  <span className="mb-1">Ankhora</span>
                  <span className="font-normal">❌ No Key</span>
                </div>
              </div>
            </div>
          )}

          {step === 3 && (
            <div className="animate-fadeInUp">
              <div className="w-24 h-24 mb-8 rounded-3xl bg-gradient-to-br from-[#C9A44A]/50 via-[#D2AF56]/60 to-[#B8934A]/70 shadow-2xl flex items-center justify-center text-5xl animate-pulse drop-shadow-lg backdrop-blur-sm">
                ✅
              </div>
              <h2 className="text-2xl md:text-3xl font-bold text-[#C9A44A] mb-4 leading-tight">
                Actions Are Cryptographically Verified
              </h2>
              <p className="text-lg md:text-xl text-zinc-700 dark:text-zinc-100/90 font-light mb-8 max-w-lg mx-auto leading-relaxed">
                Every file you upload, share, or modify is logged on the <span className="font-semibold">Stellar blockchain</span>. You can independently verify your data was encrypted and never exposed.
              </p>
              <div className="flex justify-center gap-3 flex-wrap max-w-lg mx-auto">
                <div className="bg-[#C9A44A]/20 backdrop-blur-sm rounded-2xl px-6 py-3 shadow-md text-[#C9A44A] font-semibold">
                  Upload File
                </div>
                <span className="self-center text-[#D2AF56] font-bold text-2xl">→</span>
                <div className="bg-[#C9A44A]/20 backdrop-blur-sm rounded-2xl px-6 py-3 shadow-md text-[#C9A44A] font-semibold">
                  Encrypt Locally
                </div>
                <span className="self-center text-[#D2AF56] font-bold text-2xl">→</span>
                <div className="bg-[#C9A44A]/20 backdrop-blur-sm rounded-2xl px-6 py-3 shadow-md text-[#C9A44A] font-semibold">
                  Log on Stellar
                </div>
                <span className="self-center text-[#D2AF56] font-bold text-2xl">→</span>
                <div className="bg-gradient-to-br from-green-400 to-green-300 text-black rounded-2xl px-6 py-3 font-bold shadow-lg">
                  ✅ Verifiable
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="border-t border-white/20 dark:border-zinc-700/20 pt-10 flex flex-col items-center">
          <div className="flex gap-4 mb-8">
            {[1, 2, 3].map((i) => (
              <span
                key={i}
                className={`w-4 h-4 rounded-full transition-transform duration-300 shadow-lg
                ${i === step
                    ? 'bg-gradient-to-r from-[#C9A44A] to-[#B8934A] shadow-[#C9A44A]/50 scale-125'
                    : i < step
                      ? 'bg-[#C9A44A] bg-opacity-60'
                      : 'bg-[#C9A44A] bg-opacity-20'
                  }`}
              />
            ))}
          </div>
          <div className="flex gap-6 w-full max-w-xs justify-center">
            {step > 1 && (
              <button
                onClick={prevStep}
                className="bg-[#C9A44A]/20 backdrop-blur-sm rounded-xl px-6 py-3 font-semibold text-[#C9A44A] hover:bg-opacity-40 shadow-md hover:shadow-xl transition"
              >
                Back
              </button>
            )}
            <button
              onClick={nextStep}
              className="bg-gradient-to-br from-[#C9A44A] via-[#D2AF56] to-[#B8934A] rounded-xl px-8 py-3 font-semibold text-black shadow-xl hover:shadow-2xl transition transform hover:-translate-y-1"
            >
              {step === totalSteps ? 'Start Using Ankhora' : 'Next'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );

}

export default OnboardingWizard;
