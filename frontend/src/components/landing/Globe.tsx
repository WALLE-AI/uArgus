"use client";

import { useEffect, useRef } from "react";
import createGlobe from "cobe";

export default function Globe() {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    let phi = 0;
    let width = 0;

    const onResize = () => {
      if (canvasRef.current) {
        width = canvasRef.current.offsetWidth;
      }
    };
    window.addEventListener("resize", onResize);
    onResize();

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const globe = createGlobe(canvasRef.current!, {
      devicePixelRatio: 2,
      width: width * 2,
      height: width * 2,
      phi: 0,
      theta: 0.3,
      dark: 1,
      diffuse: 3,
      mapSamples: 16000,
      mapBrightness: 1.2,
      baseColor: [0.15, 0.18, 0.3],
      markerColor: [0.231, 0.51, 0.965],
      glowColor: [0.1, 0.15, 0.3],
      markers: [
        { location: [39.9, 116.4], size: 0.06 },
        { location: [37.7, -122.4], size: 0.06 },
        { location: [51.5, -0.12], size: 0.06 },
        { location: [35.7, 139.7], size: 0.06 },
        { location: [48.9, 2.35], size: 0.04 },
        { location: [-33.9, 151.2], size: 0.04 },
        { location: [1.35, 103.8], size: 0.04 },
        { location: [55.75, 37.62], size: 0.04 },
        { location: [19.4, -99.1], size: 0.03 },
        { location: [-23.55, -46.63], size: 0.03 },
      ],
      onRender: (state: Record<string, number>) => {
        state.phi = phi;
        phi += 0.003;
        state.width = width * 2;
        state.height = width * 2;
      },
    } as Parameters<typeof createGlobe>[1]);

    return () => {
      globe.destroy();
      window.removeEventListener("resize", onResize);
    };
  }, []);

  return (
    <div className="relative mx-auto aspect-square w-full max-w-[580px]">
      <canvas
        ref={canvasRef}
        className="h-full w-full"
        style={{ contain: "layout paint size" }}
      />
      <div className="pointer-events-none absolute inset-0 rounded-full bg-gradient-to-b from-transparent via-transparent to-hero-from/80" />
    </div>
  );
}
