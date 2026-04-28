"use client";

import {
  ResponsiveContainer,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
} from "recharts";

interface AreaChartCardProps {
  title: string;
  subtitle?: string;
  data: Record<string, unknown>[];
  xKey: string;
  areas: { key: string; color: string; name?: string }[];
  height?: number;
  stacked?: boolean;
}

export default function AreaChartCard({
  title,
  subtitle,
  data,
  xKey,
  areas,
  height = 240,
  stacked = false,
}: AreaChartCardProps) {
  return (
    <div className="rounded-xl border border-border bg-white p-6">
      <div className="mb-4">
        <h3 className="text-sm font-semibold">{title}</h3>
        {subtitle && (
          <p className="text-xs text-muted-foreground mt-0.5">{subtitle}</p>
        )}
      </div>
      <ResponsiveContainer width="100%" height={height}>
        <AreaChart data={data} margin={{ top: 5, right: 5, left: -20, bottom: 0 }}>
          <defs>
            {areas.map((a) => (
              <linearGradient key={a.key} id={`grad-${a.key}`} x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor={a.color} stopOpacity={0.3} />
                <stop offset="95%" stopColor={a.color} stopOpacity={0.02} />
              </linearGradient>
            ))}
          </defs>
          <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--border))" />
          <XAxis
            dataKey={xKey}
            tick={{ fontSize: 11, fill: "hsl(var(--muted-foreground))" }}
            axisLine={false}
            tickLine={false}
          />
          <YAxis
            tick={{ fontSize: 11, fill: "hsl(var(--muted-foreground))" }}
            axisLine={false}
            tickLine={false}
          />
          <Tooltip
            contentStyle={{
              background: "white",
              border: "1px solid hsl(var(--border))",
              borderRadius: 8,
              fontSize: 12,
            }}
          />
          {areas.map((a) => (
            <Area
              key={a.key}
              type="monotone"
              dataKey={a.key}
              name={a.name ?? a.key}
              stroke={a.color}
              fill={`url(#grad-${a.key})`}
              strokeWidth={2}
              stackId={stacked ? "1" : undefined}
            />
          ))}
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
}
