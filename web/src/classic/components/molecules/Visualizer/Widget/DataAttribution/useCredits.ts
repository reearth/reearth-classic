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
};

export default ({ credits, property }: Props) => {
  const [cesiumCredit, setCesiumCredit] = useState<ProcessedCredit | undefined>(undefined);
  const [otherCredits, setOtherCredits] = useState<ProcessedCredit[] | undefined>(undefined);

  useEffect(() => {
    if (credits) {
      const cesium = processCredit(credits.engine?.cesium);
      setCesiumCredit(cesium);

      const lightboxCredits = credits.lightbox
        .map(processCredit)
        .filter(Boolean) as ProcessedCredit[];
      const screenCredits = credits.screen.map(processCredit).filter(Boolean) as ProcessedCredit[];
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
