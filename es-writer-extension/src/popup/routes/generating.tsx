import React from 'react';

const Generating: React.FC = () => {
  return (
    <div className="flex items-center justify-center p-2 rounded-md shadow-sm w-40 h-20">
      <span className="text-sm font-medium mr-2">回答生成中</span>
      <div className="flex space-x-1">
        {[0, 1, 2].map((index) => (
          <div
            key={index}
            className={`w-1.5 h-1.5 bg-blue-500 rounded-full animate-bounce`}
            style={{ animationDelay: `${index * 0.2}s` }}
          ></div>
        ))}
      </div>
    </div>
  );
};

export default Generating;