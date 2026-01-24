import React from "react";
import { FONT_FAMILY } from "../fonts";

interface TerminalWindowProps {
  children: React.ReactNode;
}

export const TerminalWindow: React.FC<TerminalWindowProps> = ({ children }) => {
  return (
    <div
      style={{
        borderRadius: 16,
        overflow: "hidden",
        boxShadow: "0 8px 32px rgba(0, 0, 0, 0.3)",
        background: "#1a1a2e",
        position: "relative",
      }}
    >

      {/* Title bar */}
      <div
        style={{
          height: 48,
          background: "#16162a",
          borderBottom: "1px solid #2a2a4a",
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
      width: 12,
      height: 12,
      borderRadius: "50%",
      backgroundColor: color,
    }}
  />
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
