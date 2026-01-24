import React, { useEffect } from "react";
import {
  AbsoluteFill,
  useCurrentFrame,
  useVideoConfig,
  interpolate,
  spring,
  delayRender,
  continueRender,
} from "remotion";
import { TerminalWindow } from "./components/TerminalWindow";
import { TypewriterText } from "./components/TypewriterText";
import { OutputDisplay } from "./components/OutputDisplay";
import { loadFonts, FONT_FAMILY } from "./fonts";

// Beeper color palette
const COLORS = {
  white: "#ffffff",
  primaryBlue: "#0c4ffb",
  cyan: "#67d9fc",
  purple: "#d074f8",
  yellow: "#d4c732",
  blue: "#4970f9",
  gradientStart: "#695ef3",
  gradientEnd: "#1b5af9",
  text: "#ffffff",
  textMuted: "rgba(255, 255, 255, 0.7)",
};

// Demo sequences showing Beeper CLI features
const SEQUENCES = [
  {
    command: "beeper chats list",
    output: "chats",
    description: "List All Chats",
    icon: "[ ]",
  },
  {
    command: "beeper messages list --chat-id abc123",
    output: "messages",
    description: "View Messages",
    icon: ">_",
  },
  {
    command: "beeper search --query \"meeting\"",
    output: "search",
    description: "Search Messages",
    icon: "?",
  },
  {
    command: "beeper send --chat-id abc123 --message \"Hello!\"",
    output: "send",
    description: "Send Messages",
    icon: ">>",
  },
];

export const BeeperCliDemo: React.FC = () => {
  const frame = useCurrentFrame();
  const { fps, durationInFrames } = useVideoConfig();

  // Load fonts
  const [handle] = React.useState(() => delayRender());
  useEffect(() => {
    loadFonts().then(() => continueRender(handle));
  }, [handle]);

  // Each sequence lasts about 3.75 seconds (112 frames at 30fps)
  const framesPerSequence = Math.floor(durationInFrames / SEQUENCES.length);
  const currentSequenceIndex = Math.floor(frame / framesPerSequence) % SEQUENCES.length;
  const frameInSequence = frame % framesPerSequence;

  const currentSequence = SEQUENCES[currentSequenceIndex];

  // Animation timing
  const typewriterDuration = 30;
  const outputDelay = 35;
  const outputFadeIn = 12;

  return (
    <AbsoluteFill
      style={{
        background: COLORS.gradientStart,
        fontFamily: FONT_FAMILY,
      }}
    >

      {/* Title */}
      <div
        style={{
          position: "absolute",
          top: 40,
          left: 0,
          right: 0,
          textAlign: "center",
          zIndex: 10,
        }}
      >
        <h1
          style={{
            fontSize: 56,
            fontWeight: 700,
            color: COLORS.white,
            margin: 0,
            letterSpacing: "-0.02em",
          }}
        >
          Beeper API CLI
        </h1>
        <p
          style={{
            fontSize: 20,
            color: COLORS.textMuted,
            marginTop: 8,
            fontWeight: 400,
          }}
        >
          Command-line interface for Beeper Desktop API
        </p>
      </div>

      {/* Terminal Window */}
      <div
        style={{
          position: "absolute",
          top: 160,
          left: "50%",
          transform: "translateX(-50%)",
          width: 1760,
        }}
      >
        <TerminalWindow>
          {/* Command Line */}
          <div style={{ marginBottom: 28, fontSize: 28 }}>
            <span style={{ color: COLORS.cyan, marginRight: 14, fontWeight: 600 }}>$</span>
            <TypewriterText
              text={currentSequence.command}
              startFrame={3}
              duration={typewriterDuration}
              frameInSequence={frameInSequence}
              color={COLORS.cyan}
            />
            <Cursor
              visible={frameInSequence >= 3 && frameInSequence < typewriterDuration + 3}
              color={COLORS.cyan}
            />
          </div>

          {/* Output */}
          {frameInSequence >= outputDelay && (
            <OutputDisplay
              type={currentSequence.output as "chats" | "messages" | "search" | "send"}
              fadeInProgress={interpolate(
                frameInSequence,
                [outputDelay, outputDelay + outputFadeIn],
                [0, 1],
                { extrapolateLeft: "clamp", extrapolateRight: "clamp" }
              )}
            />
          )}
        </TerminalWindow>
      </div>

      {/* Feature indicator */}
      <div
        style={{
          position: "absolute",
          bottom: 40,
          left: "50%",
          transform: "translateX(-50%)",
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          gap: 20,
        }}
      >
        {/* Feature label */}
        <div
          style={{
            fontSize: 24,
            color: COLORS.white,
            fontWeight: 600,
            letterSpacing: "0.08em",
            textTransform: "uppercase",
          }}
        >
          {currentSequence.description}
        </div>

        {/* Progress dots */}
        <div style={{ display: "flex", gap: 16 }}>
          {SEQUENCES.map((seq, index) => {
            const isActive = index === currentSequenceIndex;
            const dotScale = isActive
              ? spring({
                  frame: frameInSequence,
                  fps,
                  config: { damping: 20, stiffness: 200, mass: 0.3 },
                })
              : 1;

            return (
              <div
                key={index}
                style={{
                  width: isActive ? 48 : 14,
                  height: 14,
                  borderRadius: 7,
                  backgroundColor: isActive ? COLORS.cyan : COLORS.white,
                  opacity: isActive ? 1 : 0.4,
                  transform: `scale(${dotScale})`,
                  transition: "width 0.3s ease",
                  boxShadow: "none",
                }}
              />
            );
          })}
        </div>
      </div>
    </AbsoluteFill>
  );
};

// Blinking cursor
interface CursorProps {
  visible: boolean;
  color: string;
}

const Cursor: React.FC<CursorProps> = ({ visible, color }) => {
  const frame = useCurrentFrame();
  const blinkVisible = Math.floor(frame / 15) % 2 === 0;

  if (!visible) return null;

  return (
    <span
      style={{
        display: "inline-block",
        width: 3,
        height: 28,
        backgroundColor: color,
        marginLeft: 4,
        opacity: blinkVisible ? 1 : 0,
        verticalAlign: "middle",
      }}
    />
  );
};
