"use client";

import {
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
} from "recharts";

interface BarChartCardProps {
  title: string;
  subtitle?: string;
  data: Record<string, unknown>[];
  xKey: string;
  bars: { key: string; color: string; name?: string }[];
  height?: number;
  layout?: "horizontal" | "vertical";
  stacked?: boolean;
}

export default function BarChartCard({
  title,
  subtitle,
  data,
  xKey,
  bars,
  height = 240,
  layout = "horizontal",
  stacked = false,
}: BarChartCardProps) {
  return (
    <div className="rounded-xl border border-border bg-white p-6">
      <div className="mb-4">
        <h3 className="text-sm font-semibold">{title}</h3>
        {subtitle && (
          <p className="text-xs text-muted-foreground mt-0.5">{subtitle}</p>
        )}
      </div>
      <ResponsiveContainer width="100%" height={height}>
        <BarChart
          data={data}
          layout={layout === "vertical" ? "vertical" : "horizontal"}
          margin={{ top: 5, right: 5, left: layout === "vertical" ? 60 : -20, bottom: 0 }}
        >
          <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--border))" />
          {layout === "vertical" ? (
            <>
              <XAxis type="number" tick={{ fontSize: 11, fill: "hsl(var(--muted-foreground))" }} axisLine={false} tickLine={false} />
              <YAxis dataKey={xKey} type="category" tick={{ fontSize: 11, fill: "hsl(var(--muted-foreground))" }} axisLine={false} tickLine={false} width={80} />
            </>
          ) : (
            <>
              <XAxis dataKey={xKey} tick={{ fontSize: 11, fill: "hsl(var(--muted-foreground))" }} axisLine={false} tickLine={false} />
              <YAxis tick={{ fontSize: 11, fill: "hsl(var(--muted-foreground))" }} axisLine={false} tickLine={false} />
            </>
          )}
          <Tooltip
            contentStyle={{
              background: "white",
              border: "1px solid hsl(var(--border))",
              borderRadius: 8,
              fontSize: 12,
            }}
          />
          {bars.map((b) => (
            <Bar
              key={b.key}
              dataKey={b.key}
              name={b.name ?? b.key}
              fill={b.color}
              radius={[4, 4, 0, 0]}
              stackId={stacked ? "1" : undefined}
            />
          ))}
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
