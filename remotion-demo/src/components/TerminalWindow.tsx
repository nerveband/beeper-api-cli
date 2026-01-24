import React from "react";
import { useCurrentFrame, interpolate } from "remotion";
import { FONT_FAMILY } from "../fonts";

interface TerminalWindowProps {
  children: React.ReactNode;
}

export const TerminalWindow: React.FC<TerminalWindowProps> = ({ children }) => {
  const frame = useCurrentFrame();

  // Subtle floating animation
  const floatY = interpolate(Math.sin(frame / 40), [-1, 1], [-4, 4]);
  const floatRotate = interpolate(Math.sin(frame / 50), [-1, 1], [-0.3, 0.3]);

  return (
    <div
      style={{
        transform: `translateY(${floatY}px) rotate(${floatRotate}deg)`,
        borderRadius: 20,
        overflow: "hidden",
        boxShadow: `
          0 0 0 1px rgba(255, 255, 255, 0.18),
          0 0 0 2px rgba(255, 255, 255, 0.05),
          0 4px 12px rgba(0, 0, 0, 0.15),
          0 16px 48px rgba(0, 0, 0, 0.2),
          0 32px 80px rgba(0, 0, 0, 0.15),
          inset 0 1px 1px rgba(255, 255, 255, 0.15),
          inset 0 -1px 1px rgba(0, 0, 0, 0.1)
        `,
        background: `
          linear-gradient(165deg,
            rgba(20, 25, 45, 0.95) 0%,
            rgba(15, 20, 40, 0.98) 30%,
            rgba(10, 15, 35, 0.98) 70%,
            rgba(8, 12, 30, 0.95) 100%
          )
        `,
        backdropFilter: "blur(50px) saturate(180%)",
        WebkitBackdropFilter: "blur(50px) saturate(180%)",
        position: "relative",
      }}
    >
      {/* Glass reflection */}
      <div
        style={{
          position: "absolute",
          top: 0,
          left: 0,
          right: 0,
          height: "50%",
          background: "linear-gradient(180deg, rgba(255,255,255,0.04) 0%, transparent 100%)",
          borderRadius: "20px 20px 0 0",
          pointerEvents: "none",
        }}
      />

      {/* Title bar */}
      <div
        style={{
          height: 52,
          background: `
            linear-gradient(180deg,
              rgba(40, 50, 80, 0.95) 0%,
              rgba(30, 40, 70, 0.95) 50%,
              rgba(25, 35, 65, 0.95) 100%
            )
          `,
          borderBottom: "1px solid rgba(255, 255, 255, 0.08)",
          display: "flex",
          alignItems: "center",
          padding: "0 20px",
          position: "relative",
        }}
      >
        {/* Traffic lights */}
        <div style={{ display: "flex", gap: 8, zIndex: 1 }}>
          <TrafficLight color="#ff5f57" />
          <TrafficLight color="#febc2e" />
          <TrafficLight color="#28c840" />
        </div>

        {/* Title */}
        <div
          style={{
            position: "absolute",
            left: "50%",
            transform: "translateX(-50%)",
            color: "rgba(255, 255, 255, 0.5)",
            fontSize: 13,
            fontWeight: 500,
            letterSpacing: "0.03em",
            fontFamily: FONT_FAMILY,
          }}
        >
          beeper — Terminal
        </div>

        {/* Window buttons */}
        <div style={{ marginLeft: "auto", display: "flex", gap: 16 }}>
          <WindowButton icon="−" />
          <WindowButton icon="□" />
        </div>
      </div>

      {/* Terminal content */}
      <div
        style={{
          padding: "40px 50px",
          fontSize: 26,
          lineHeight: 1.8,
          color: "#e8e8e8",
          minHeight: 620,
          fontFamily: FONT_FAMILY,
          position: "relative",
        }}
      >
        {children}
      </div>
    </div>
  );
};

interface TrafficLightProps {
  color: string;
}

const TrafficLight: React.FC<TrafficLightProps> = ({ color }) => (
  <div
    style={{
      width: 13,
      height: 13,
      borderRadius: "50%",
      backgroundColor: color,
      boxShadow: `
        inset 0 1px 0 rgba(255, 255, 255, 0.4),
        inset 0 -1px 2px rgba(0, 0, 0, 0.25),
        0 1px 2px rgba(0, 0, 0, 0.2)
      `,
      position: "relative",
    }}
  >
    <div
      style={{
        position: "absolute",
        top: 2,
        left: 3,
        width: 5,
        height: 3,
        borderRadius: "50%",
        background: "linear-gradient(180deg, rgba(255,255,255,0.5) 0%, transparent 100%)",
      }}
    />
  </div>
);

interface WindowButtonProps {
  icon: string;
}

const WindowButton: React.FC<WindowButtonProps> = ({ icon }) => (
  <div
    style={{
      color: "rgba(255, 255, 255, 0.3)",
      fontSize: 14,
      fontWeight: 300,
    }}
  >
    {icon}
  </div>
);
