"use client";

import React, { useState } from "react";

const Grid: React.FC = () => {
  const numRows = 6;
  const numCols = 7;

  type ClickedState = {
    [key: string]: number;
  };

  const [clicked, setClicked] = useState<ClickedState>({});
  const [player, setPlayer] = useState<number>(1);
  const [winner, setWinner] = useState<number | null>(null);

  const checkWin = (row: number, col: number) => {
    const directions = [
      { r: 1, c: 0 }, // Vertical
      { r: 0, c: 1 }, // Horizontal
      { r: 1, c: 1 }, // Diagonal Down-Right
      { r: 1, c: -1 }, // Diagonal Down-Left
    ];

    const checkDirection = (dr: number, dc: number) => {
      let count = 1;

      for (let step = 1; step < 4; step++) {
        const r = row + step * dr;
        const c = col + step * dc;
        if (clicked[`${r}-${c}`] === player) {
          count++;
        } else {
          break;
        }
      }

      for (let step = 1; step < 4; step++) {
        const r = row - step * dr;
        const c = col - step * dc;
        if (clicked[`${r}-${c}`] === player) {
          count++;
        } else {
          break;
        }
      }

      return count >= 4;
    };

    for (const { r: dr, c: dc } of directions) {
      if (checkDirection(dr, dc)) {
        return true;
      }
    }

    return false;
  };

  const handleClick = (col: number) => {
    if (winner !== null) return;

    for (let row = numRows - 1; row >= 0; row--) {
      let currentKey = `${row}-${col}`;
      if (clicked[currentKey] == undefined) {
        setClicked((prev) => {
          const newState = { ...prev, [currentKey]: player };
          if (checkWin(row, col)) {
            setWinner(player);
          }
          return newState;
        });
        setPlayer((prev) => (prev === 1 ? 2 : 1));
        return;
      }
    }
  };

  const gridItems = [];
  for (let row = 0; row < numRows; row++) {
    for (let col = 0; col < numCols; col++) {
      const key = `${row}-${col}`;
      const cellClass =
        clicked[key] === 0
          ? "bg-white"
          : clicked[key] === 1
          ? "bg-red-500"
          : clicked[key] === 2
          ? "bg-blue-500"
          : "bg-white";

      gridItems.push(
        <div
          key={key}
          onClick={() => handleClick(col)}
          className='w-32 h-32 flex items-center justify-center'
        >
          <div className={`w-3/4 h-3/4 rounded-full ${cellClass}`}></div>
        </div>
      );
    }
  }

  return (
    <div>
      {winner !== null && (
        <div className='text-center text-2xl font-bold mb-4'>
          Player {winner} Wins!
        </div>
      )}
      <div className='bg-blue-400 grid grid-cols-7'>{gridItems}</div>
    </div>
  );
};

const Page: React.FC = () => {
  return (
    <main className='grid justify-center h-screen items-center'>
      <Grid />
    </main>
  );
};

export default Page;
