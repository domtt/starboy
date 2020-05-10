import React, { useState, useEffect } from 'react';
import './App.css';

import { LineChart, XAxis, YAxis, Line, Legend } from "recharts";
import { useToken } from './state/token';

const getStars = (repo: string, token: string) =>
  fetch(`http://localhost:8080/api/repos/${repo}/stars?token=${token}`, {
    headers: {
      Authorization: token
    }
  });

interface StarHistoryEntry {
  unixTime: number;
  stars: number;
}

const App = () => {
  const token = useToken();

  const [newRepo, setNewRepo] = useState("");
  const [data, setData] = useState<{ [k: string]: StarHistoryEntry[] }>({});
  const mergeData = (o: any) => setData({...data, ...o});
  const fetchStars = () => {
    if (token === null) {
      return;
    }
    getStars(newRepo, token)
      .then(r => r.json())
      .then(mergeData);
    setNewRepo("");
  };

  useEffect(() => console.log(data), [data]);

  const all = Object.keys(data)
    .map(k => data[k])
    .reduce((a, b) => [...a, ...b], []);
  const times: number[] = all.map(d => d.unixTime);
  let [min, max] = [Math.min(...times), Math.max(...times)];
  let startYear = new Date(1000 * min).getFullYear();
  let endYear = new Date(1000 * max).getFullYear();
  let years = Array.from(
    { length: 1 + endYear - startYear },
    (_, i) => i + startYear
  );

  return (
    <div>
      <form onSubmit={e => {e.preventDefault();fetchStars()}}>
      <input type="text" value={newRepo} onChange={e => setNewRepo(e.target.value)} placeholder="repository" />
      <button disabled={!token && false}>
        Fetch stars
      </button>
      </form>
      {Object.keys(data).length > 0 && (
        <LineChart width={800} height={400} data={all}>
          <XAxis
            dataKey="unixTime"
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
              dataKey={"stars"}
              name={k}
              data={data[k].sort((a, b) => a.unixTime - b.unixTime)}
            />
          ))}
        </LineChart>
      )}
    </div>
  );
};

export default App;
