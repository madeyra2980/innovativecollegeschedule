import React from 'react';
import TrelloSchedule from './TrelloSchedule';

const Schedule = () => {
  return (
    <div>
      <div className="card">
        <div className="schedule-header">
          <h2>Расписание</h2>
        </div>

        <TrelloSchedule />
      </div>
    </div>
  );
};

export default Schedule;