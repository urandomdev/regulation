import React from 'react';
import { useNavigate } from 'react-router-dom';
import './UpcomingMissions.css';

const UpcomingMissions = () => {
  const navigate = useNavigate();

  const missions = [
    {
      id: 1,
      title: 'Website Redesign',
      client: 'TechCorp Solutions',
      due: 'Nov 30, 2025'
    },
    {
      id: 2,
      title: 'Mobile App Mockups',
      client: 'StartupHub Inc',
      due: 'Dec 5, 2025'
    }
  ];

  const handleMissionClick = () => {
    navigate('/saving-goals');
  };

  const handleViewAll = () => {
    navigate('/saving-goals');
  };

  return (
    <div className="upcoming-missions-container">
      <div className="missions-header">
        <h2>Upcoming Missions</h2>
        <button className="view-all-link" onClick={handleViewAll}>
          View All
        </button>
      </div>

      <div className="missions-list">
        {missions.map((mission) => (
          <div key={mission.id} className="mission-item" onClick={handleMissionClick}>
            <div className="mission-info">
              <h3 className="mission-title">{mission.title}</h3>
              <p className="mission-client">Client: {mission.client}</p>
              <p className="mission-due">Due: {mission.due}</p>
            </div>
            <svg className="mission-chevron" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
            </svg>
          </div>
        ))}
      </div>
    </div>
  );
};

export default UpcomingMissions;