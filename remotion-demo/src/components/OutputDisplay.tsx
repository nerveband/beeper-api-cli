import React from "react";

interface OutputDisplayProps {
  type: "chats" | "messages" | "search" | "send";
  fadeInProgress: number;
}

// Beeper colors
const COLORS = {
  cyan: "#67d9fc",
  purple: "#d074f8",
  yellow: "#d4c732",
  blue: "#4970f9",
  string: "#a8e6cf",
  key: "#67d9fc",
  bracket: "#d074f8",
  network: "#d4c732",
  success: "#28c840",
};

export const OutputDisplay: React.FC<OutputDisplayProps> = ({
  type,
  fadeInProgress,
}) => {
  return (
    <div
      style={{
        opacity: fadeInProgress,
        transform: `translateY(${(1 - fadeInProgress) * 10}px)`,
      }}
    >
      {type === "chats" && <ChatsOutput />}
      {type === "messages" && <MessagesOutput />}
      {type === "search" && <SearchOutput />}
      {type === "send" && <SendOutput />}
    </div>
  );
};

const ChatsOutput: React.FC = () => (
  <pre
    style={{
      margin: 0,
      fontFamily: "inherit",
      fontSize: 22,
      lineHeight: 1.6,
    }}
  >
    <span style={{ color: COLORS.bracket }}>{"["}</span>
    {"\n  "}
    <span style={{ color: COLORS.bracket }}>{"{"}</span>
    {"\n    "}
    <span style={{ color: COLORS.key }}>"id"</span>
    <span style={{ color: "#e0e0e0" }}>: </span>
    <span style={{ color: COLORS.string }}>"!abc123:beeper.local"</span>
    <span style={{ color: "#e0e0e0" }}>,</span>
    {"\n    "}
    <span style={{ color: COLORS.key }}>"title"</span>
    <span style={{ color: "#e0e0e0" }}>: </span>
    <span style={{ color: COLORS.string }}>"Team Chat"</span>
    <span style={{ color: "#e0e0e0" }}>,</span>
    {"\n    "}
    <span style={{ color: COLORS.key }}>"network"</span>
    <span style={{ color: "#e0e0e0" }}>: </span>
    <span style={{ color: COLORS.network }}>"Telegram"</span>
    <span style={{ color: "#e0e0e0" }}>,</span>
    {"\n    "}
    <span style={{ color: COLORS.key }}>"unreadCount"</span>
    <span style={{ color: "#e0e0e0" }}>: </span>
    <span style={{ color: COLORS.yellow }}>3</span>
    {"\n  "}
    <span style={{ color: COLORS.bracket }}>{"}"}</span>
    <span style={{ color: "#e0e0e0" }}>,</span>
    {"\n  "}
    <span style={{ color: COLORS.bracket }}>{"{"}</span>
    {"\n    "}
    <span style={{ color: COLORS.key }}>"id"</span>
    <span style={{ color: "#e0e0e0" }}>: </span>
    <span style={{ color: COLORS.string }}>"!def456:beeper.local"</span>
    <span style={{ color: "#e0e0e0" }}>,</span>
    {"\n    "}
    <span style={{ color: COLORS.key }}>"title"</span>
    <span style={{ color: "#e0e0e0" }}>: </span>
    <span style={{ color: COLORS.string }}>"Family Group"</span>
    <span style={{ color: "#e0e0e0" }}>,</span>
    {"\n    "}
    <span style={{ color: COLORS.key }}>"network"</span>
    <span style={{ color: "#e0e0e0" }}>: </span>
    <span style={{ color: COLORS.network }}>"WhatsApp"</span>
    <span style={{ color: "#e0e0e0" }}>,</span>
    {"\n    "}
    <span style={{ color: COLORS.key }}>"unreadCount"</span>
    <span style={{ color: "#e0e0e0" }}>: </span>
    <span style={{ color: COLORS.yellow }}>0</span>
    {"\n  "}
    <span style={{ color: COLORS.bracket }}>{"}"}</span>
    {"\n"}
    <span style={{ color: COLORS.bracket }}>{"]"}</span>
  </pre>
);

const MessagesOutput: React.FC = () => (
  <pre
    style={{
      margin: 0,
      fontFamily: "inherit",
      fontSize: 22,
      lineHeight: 1.6,
    }}
  >
    <span style={{ color: COLORS.bracket }}>{"["}</span>
    {"\n  "}
    <span style={{ color: COLORS.bracket }}>{"{"}</span>
    {"\n    "}
    <span style={{ color: COLORS.key }}>"sender"</span>
    <span style={{ color: "#e0e0e0" }}>: </span>
    <span style={{ color: COLORS.string }}>"Alice"</span>
    <span style={{ color: "#e0e0e0" }}>,</span>
    {"\n    "}
    <span style={{ color: COLORS.key }}>"text"</span>
    <span style={{ color: "#e0e0e0" }}>: </span>
    <span style={{ color: COLORS.string }}>"Hey, are you coming to the meeting?"</span>
    <span style={{ color: "#e0e0e0" }}>,</span>
    {"\n    "}
    <span style={{ color: COLORS.key }}>"timestamp"</span>
    <span style={{ color: "#e0e0e0" }}>: </span>
    <span style={{ color: COLORS.purple }}>"2024-01-15T10:30:00Z"</span>
    {"\n  "}
    <span style={{ color: COLORS.bracket }}>{"}"}</span>
    <span style={{ color: "#e0e0e0" }}>,</span>
    {"\n  "}
    <span style={{ color: COLORS.bracket }}>{"{"}</span>
    {"\n    "}
    <span style={{ color: COLORS.key }}>"sender"</span>
    <span style={{ color: "#e0e0e0" }}>: </span>
    <span style={{ color: COLORS.string }}>"You"</span>
    <span style={{ color: "#e0e0e0" }}>,</span>
    {"\n    "}
    <span style={{ color: COLORS.key }}>"text"</span>
    <span style={{ color: "#e0e0e0" }}>: </span>
    <span style={{ color: COLORS.string }}>"Yes, on my way!"</span>
    <span style={{ color: "#e0e0e0" }}>,</span>
    {"\n    "}
    <span style={{ color: COLORS.key }}>"timestamp"</span>
    <span style={{ color: "#e0e0e0" }}>: </span>
    <span style={{ color: COLORS.purple }}>"2024-01-15T10:31:00Z"</span>
    {"\n  "}
    <span style={{ color: COLORS.bracket }}>{"}"}</span>
    {"\n"}
    <span style={{ color: COLORS.bracket }}>{"]"}</span>
  </pre>
);

const SearchOutput: React.FC = () => (
  <div style={{ fontFamily: "inherit", fontSize: 22 }}>
    <div style={{ color: "rgba(255,255,255,0.5)", marginBottom: 20 }}>
      Found <span style={{ color: COLORS.cyan }}>3</span> results for "meeting"
    </div>
    <div style={{ marginBottom: 16, paddingLeft: 12, borderLeft: `3px solid ${COLORS.purple}` }}>
      <div style={{ color: COLORS.network, fontSize: 14, marginBottom: 4 }}>Telegram · Team Chat</div>
      <div style={{ color: "#e0e0e0" }}>
        "Hey, are you coming to the <span style={{ color: COLORS.yellow, fontWeight: 600 }}>meeting</span>?"
      </div>
    </div>
    <div style={{ marginBottom: 16, paddingLeft: 12, borderLeft: `3px solid ${COLORS.cyan}` }}>
      <div style={{ color: COLORS.network, fontSize: 14, marginBottom: 4 }}>WhatsApp · Family Group</div>
      <div style={{ color: "#e0e0e0" }}>
        "Don't forget the family <span style={{ color: COLORS.yellow, fontWeight: 600 }}>meeting</span> tomorrow"
      </div>
    </div>
    <div style={{ paddingLeft: 12, borderLeft: `3px solid ${COLORS.blue}` }}>
      <div style={{ color: COLORS.network, fontSize: 14, marginBottom: 4 }}>iMessage · Work</div>
      <div style={{ color: "#e0e0e0" }}>
        "<span style={{ color: COLORS.yellow, fontWeight: 600 }}>Meeting</span> notes attached"
      </div>
    </div>
  </div>
);

const SendOutput: React.FC = () => (
  <div style={{ fontFamily: "inherit", fontSize: 22 }}>
    <div
      style={{
        display: "flex",
        alignItems: "center",
        gap: 12,
        marginBottom: 20,
        color: COLORS.success,
      }}
    >
      <span style={{ fontSize: 28 }}>✓</span>
      <span>Message sent successfully</span>
    </div>
    <pre
      style={{
        margin: 0,
        fontFamily: "inherit",
        fontSize: 20,
        lineHeight: 1.6,
      }}
    >
      <span style={{ color: COLORS.bracket }}>{"{"}</span>
      {"\n  "}
      <span style={{ color: COLORS.key }}>"id"</span>
      <span style={{ color: "#e0e0e0" }}>: </span>
      <span style={{ color: COLORS.string }}>"msg_12345"</span>
      <span style={{ color: "#e0e0e0" }}>,</span>
      {"\n  "}
      <span style={{ color: COLORS.key }}>"chatId"</span>
      <span style={{ color: "#e0e0e0" }}>: </span>
      <span style={{ color: COLORS.string }}>"!abc123:beeper.local"</span>
      <span style={{ color: "#e0e0e0" }}>,</span>
      {"\n  "}
      <span style={{ color: COLORS.key }}>"text"</span>
      <span style={{ color: "#e0e0e0" }}>: </span>
      <span style={{ color: COLORS.string }}>"Hello!"</span>
      <span style={{ color: "#e0e0e0" }}>,</span>
      {"\n  "}
      <span style={{ color: COLORS.key }}>"status"</span>
      <span style={{ color: "#e0e0e0" }}>: </span>
      <span style={{ color: COLORS.success }}>"delivered"</span>
      {"\n"}
      <span style={{ color: COLORS.bracket }}>{"}"}</span>
    </pre>
  </div>
);
