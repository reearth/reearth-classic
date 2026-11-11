import { useEffect, useState } from "react";

import { Credits, CreditItem } from "../../Engine/ref";

import { DataAttributionWidgetProperty } from ".";

type Props = {
  credits?: Credits;
  property?: DataAttributionWidgetProperty;
};

export type ProcessedCredit = {
  logo?: string;
  description?: string;
  creditUrl?: string;
  builtinHtml?: string; // For hardcoded credits like Re:Earth Buildings
};

export default ({ credits, property }: Props) => {
  const [cesiumCredit, setCesiumCredit] = useState<ProcessedCredit | undefined>(undefined);
  const [otherCredits, setOtherCredits] = useState<ProcessedCredit[] | undefined>(undefined);
  const [googleCredit, setGoogleCredit] = useState<ProcessedCredit | undefined>(undefined);

  useEffect(() => {
    if (credits) {
      const cesium = processCredit(credits.engine?.cesium);
      setCesiumCredit(cesium);

      const lightboxCredits = credits.lightbox
        .map(processCredit)
        .filter(Boolean) as ProcessedCredit[];
      const screenCredits = credits.screen.map(processCredit).filter(Boolean) as ProcessedCredit[];
      let googleCredit = screenCredits.find(e => e.logo?.includes("google"));
      if (!googleCredit && credits.lightbox.some(c => c.html?.includes("Google"))) {
        googleCredit = {
          logo: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAGIAAAASCAYAAACghwvPAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAs1JREFUeNrsWYGRhCAMtARLsARLoARLoARLsANKsARLoARKsARK4M+fcLOXCxHP+fPm/5lh5tQAIZtkA9c00FJK5taX9Ny2d0Nzcbvp0N+633rzW9ttcy7tt/liHU1W5MQc09Z3ZFqSs1eCEDbv35ShbwN54W8BIjeryMwk4y/ZnGbot3vHzwMRs7MJ6S9dAUTmhFVS7KCRukq5dkem2+SOAEFj+gNAJClFQfS/HYjcxhfGthDGCQAdBCMtgpwpkDF67ULvnAQEpc6VjRl39htgTMfmSqCDV9Z52CfoOJEOWReHjkfPEbPQ93d4YQreIXVHIASF2C2AFZmhsA0AQlTm8xwIMBwaR+UAmMvmihC+5fE9AsHW8cxZjBBpgWQirgFcvNLvcF9nBwjNKGPB6Bgh+BxBaQRxZSkhZC9l4EhArPCtZam2xAFo4KyDIU++8yST62kfFjmTGfmJZ9k+sejpWIRYNTWRgthnwStmIV1l41n4PQrpKoEhHiKElZtPQDBSNULIB4mzmIENyEYEr5CaFtj7wub6Bl9YL0eTA/vliLhXp4fIGjxwAiA0wpsqI25SSNgUgDCls45G2oKBF4m8mZGtkJpWAQi/o7/EqTFHhHkijvJBCL1Ri4gkRMSkRER/MiJWSpXZm9sDQHSYIgtAeM47sM/aiJgFgNzDOHagW8mAHS02MHJyNGbkaY3kPeZoxhEDbL7EEcglppIjkFeMVK2UgAAO6CqAaEuHPm4fgSMsRBLOEx4yAfP4qiuOA1WTJmdOVE1GkQ+1QOxFDnOmhVVnvlA1zehEharJS9HYACHHmks/MjInSekcIeVGz3M5gbGwuQLK80s/4eyxl2I9eq0ChIfo5/pH8HDHyuJJPCeUzxG+5iD67isMK4DjeL3/oZem7z2J/+BGZiXfp0+46/orQGgcEc7cg/0D8fofP1X5/sN0P/X/xZcAAwB8nuN904Y57wAAAABJRU5ErkJggg==",
          description: "Google Maps",
        };
      }
      setGoogleCredit(googleCredit);
      const widgetCredits =
        property?.default?.map(credit => ({
          logo: credit.logo,
          description: credit.description,
          creditUrl: credit.creditUrl,
        })) || [];

      setOtherCredits([...lightboxCredits, ...screenCredits, ...widgetCredits]);
    }
  }, [credits, property]);

  return {
    cesiumCredit,
    otherCredits,
    googleCredit,
  };
};

function processCredit(credit: CreditItem | undefined): ProcessedCredit | undefined {
  if (!credit) return undefined;

  if (credit?.html) {
    try {
      const parser = new DOMParser();
      const doc = parser.parseFromString(credit.html, "text/html");
      const logo = doc.querySelector("img")?.getAttribute("src") || "";
      const description =
        doc.querySelector("img")?.getAttribute("title") || doc.body.textContent || "";
      const creditUrl = doc.querySelector("a")?.getAttribute("href") || "";

      return { logo, description, creditUrl };
    } catch (error) {
      console.error("Error processing credit HTML:", error);
      return undefined;
    }
  }

  return undefined;
}
