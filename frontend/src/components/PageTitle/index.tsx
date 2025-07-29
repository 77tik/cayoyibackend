import React from 'react';
import styled from 'styled-components';

interface PageTitleProps {
  title: string;
  className?: string;
}

const TitleWrapper = styled.div`
  color: #fff;
  /* margin-left: 50px; */

  font-size: 20px;
  font-family: YouSheBiaoTiHei;
  font-size: 36px;
  letter-spacing: 6px;
  line-height: 36px;
  text-shadow: 2px 2px 2px rgba(10, 48, 162, 0.4);
  font-weight: 400;
`;

const PageTitle: React.FC<PageTitleProps> = ({ title, className }) => {
  return <TitleWrapper className={className}>{title}</TitleWrapper>;
};

export default PageTitle;
