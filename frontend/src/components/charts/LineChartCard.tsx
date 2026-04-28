"use client";

import {
  ResponsiveContainer,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
} from "recharts";

interface LineChartCardProps {
  title: string;
  subtitle?: string;
  data: Record<string, unknown>[];
  xKey: string;
  lines: { key: string; color: string; name?: string; dashed?: boolean }[];
  height?: number;
}

export default function LineChartCard({
  title,
  subtitle,
  data,
  xKey,
  lines,
  height = 240,
}: LineChartCardProps) {
  return (
    <div className="rounded-xl border border-border bg-white p-6">
      <div className="mb-4">
        <h3 className="text-sm font-semibold">{title}</h3>
        {subtitle && (
          <p className="text-xs text-muted-foreground mt-0.5">{subtitle}</p>
        )}
      </div>
      <ResponsiveContainer width="100%" height={height}>
        <LineChart data={data} margin={{ top: 5, right: 5, left: -20, bottom: 0 }}>
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
          {lines.length > 1 && (
            <Legend iconType="circle" iconSize={8} wrapperStyle={{ fontSize: 11 }} />
          )}
          {lines.map((l) => (
            <Line
              key={l.key}
              type="monotone"
              dataKey={l.key}
              name={l.name ?? l.key}
              stroke={l.color}
              strokeWidth={2}
              strokeDasharray={l.dashed ? "5 5" : undefined}
              dot={false}
              activeDot={{ r: 4 }}
            />
          ))}
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
