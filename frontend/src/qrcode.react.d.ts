declare module "qrcode.react" {
  import * as React from "react";
  interface QRCodeProps {
    value: string;
    size?: number;
    bgColor?: string;
    fgColor?: string;
    level?: "L" | "M" | "Q" | "H";
    renderAs?: "canvas" | "svg";
    includeMargin?: boolean;
    imageSettings?: {
      src: string;
      x?: number;
      y?: number;
      height?: number;
      width?: number;
      excavate?: boolean;
    };
    style?: React.CSSProperties;
    className?: string;
  }
  export default class QRCode extends React.Component<QRCodeProps, any> {}
}
