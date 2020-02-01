import React, { useEffect, useState } from "react";
import "./App.css";

import { LineChart, XAxis, YAxis, Line, Legend } from "recharts";

const getStars = (repos: string[], token: string) =>
  fetch(`http://localhost:3000/api/data?repos=${repos.join(",")}`, {
    headers: {
      Authorization: token
    }
  });

interface StarHistoryEntry {
  t: number;
  v: number;
}

const App = () => {
  const [repos] = useState<string[]>([
    "angular/angular",
    "facebook/react"
  ]);
  const [token, setToken] = useState<string>("");
  const [data, setData] = useState<{ [k: string]: StarHistoryEntry[] }>({});
  const fetchStars = () => {
    console.log(token)
    getStars(repos, token)
      .then(r => r.json())
      .then(setData);
  };

  useEffect(() => console.log(data), [data]);

  const all = Object.keys(data)
    .map(k => data[k])
    .reduce((a, b) => [...a, ...b], []);
  const times: number[] = all.map(d => d.t);
  let [min, max] = [Math.min(...times), Math.max(...times)];
  let startYear = new Date(1000 * min).getFullYear();
  let endYear = new Date(1000 * max).getFullYear();
  let years = Array.from(
    { length: 1 + endYear - startYear },
    (_, i) => i + startYear
  );

  return (
    <div className="App">
      <input type="text" value={token} onChange={e => setToken(e.target.value)} placeholder="token" /> 
      <button onClick={fetchStars}>Fetch stars</button>
      {Object.keys(data).length > 0 && (
        <LineChart width={800} height={400} data={all}>
          <XAxis
            dataKey="t"
            domain={["dataMin", "dataMax"]}
            name="Time"
            scale="time"
            type="number"
            interval={"preserveStartEnd"}
            ticks={years.map(y => new Date(y, 0).getTime() / 1000)}
            tickFormatter={unixTime => new Date(1000 * unixTime).getFullYear()}
          />
          <YAxis
            tickFormatter={s => (s === 0 ? "0" : Math.floor(s / 1000) + "k")}
          />
          <Legend />
          {Object.keys(data).map(k => (
            <Line
              key={k}
              dataKey={"v"}
              name={k}
              data={data[k].sort((a, b) => a.t - b.t)}
            />
          ))}
        </LineChart>
      )}
    </div>
  );
};

export default App;
