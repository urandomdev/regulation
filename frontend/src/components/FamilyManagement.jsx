import React from 'react';
import './FamilyManagement.css';

const FamilyManagement = () => {
  const familyMembers = [
    {
      id: 1,
      name: 'Security Audit',
      description: 'Review system access logs',
      status: 'attention'
    },
    {
      id: 2,
      name: 'Compliance Check',
      description: 'Q4 regulatory review',
      status: 'in-progress'
    },
    {
      id: 3,
      name: 'Data Backup',
      description: 'Weekly backup verification',
      status: 'attention'
    }
  ];

  return (
    <div className="family-management-container">
      <h2>Family Management</h2>
      <div className="family-members-list">
        {familyMembers.map((member) => (
          <div key={member.id} className="family-member-item">
            <div className="family-member-info">
              <h3 className="family-member-name">{member.name}</h3>
              <p className="family-member-description">{member.description}</p>
            </div>
            <span className={`status-badge ${member.status}`}>
              {member.status === 'attention' ? 'Requires Attention' : 'In Progress'}
            </span>
          </div>
        ))}
      </div>
    </div>
  );
};

export default FamilyManagement;