export default function TasksPage() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">采集任务</h1>
      <div className="rounded-xl border border-border bg-white">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border text-left text-xs text-muted-foreground">
              <th className="px-6 py-3 font-medium">任务名称</th>
              <th className="px-6 py-3 font-medium">频率</th>
              <th className="px-6 py-3 font-medium">状态</th>
              <th className="px-6 py-3 font-medium">上次执行</th>
              <th className="px-6 py-3 font-medium">下次执行</th>
            </tr>
          </thead>
          <tbody>
            {[
              { name: "RSS 全量抓取", freq: "每 15 分钟", status: "运行中", last: "2 min ago", next: "13 min" },
              { name: "API 数据同步", freq: "每小时", status: "等待", last: "48 min ago", next: "12 min" },
              { name: "深度爬取", freq: "每 6 小时", status: "运行中", last: "2 hr ago", next: "4 hr" },
              { name: "OPML 同步", freq: "每日", status: "等待", last: "18 hr ago", next: "6 hr" },
            ].map((t) => (
              <tr key={t.name} className="border-b border-border last:border-0">
                <td className="px-6 py-3 font-medium">{t.name}</td>
                <td className="px-6 py-3 text-muted-foreground">{t.freq}</td>
                <td className="px-6 py-3">
                  <span className={`inline-flex items-center gap-1 text-xs ${t.status === "运行中" ? "text-accent-green" : "text-muted-foreground"}`}>
                    <span className={`h-1.5 w-1.5 rounded-full ${t.status === "运行中" ? "bg-accent-green animate-pulse" : "bg-muted-foreground/40"}`} />
                    {t.status}
                  </span>
                </td>
                <td className="px-6 py-3 text-muted-foreground">{t.last}</td>
                <td className="px-6 py-3 text-muted-foreground">{t.next}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
