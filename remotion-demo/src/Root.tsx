import { Composition } from "remotion";
import { BeeperCliDemo } from "./BeeperCliDemo";

export const RemotionRoot: React.FC = () => {
  return (
    <>
      <Composition
        id="BeeperCliDemo"
        component={BeeperCliDemo}
        durationInFrames={450}
        fps={30}
        width={1920}
        height={1080}
      />
    </>
  );
};
